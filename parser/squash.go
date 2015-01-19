package parser

import (
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/winchman/builder-core/communication"
	"github.com/winchman/libsquash"
)

type squashImageOptions struct {
	Image       *docker.APIImages
	PipeWriter  *io.PipeWriter
	Opts        *DockerCmdOpts
	RetIDBuffer io.Writer
}

func (b *BuildCmd) squashImage(squashOpts squashImageOptions) error {
	// the image to be squashed
	image := squashOpts.Image

	// access to the docker api
	opts := squashOpts.Opts

	// the buffer on which to return the id of the squashed image
	retIDBuffer := squashOpts.RetIDBuffer

	// create some pipes
	imageReader, pipeWriter := io.Pipe()
	squashedImageReader, squashedImageWriter := io.Pipe()

	// defer closing all pipes
	defer func() {
		imageReader.Close()
		squashedImageReader.Close()
		pipeWriter.Close()
		squashedImageWriter.Close()
	}()

	// exporting async - error can be ignored because errors will be propagated
	// on the pipes.  all pipes are synchronous
	go b.exportImage(image, pipeWriter, opts)

	// squash async - see note above about ignoring errors
	go b.squash(imageReader, squashedImageWriter, retIDBuffer)

	// import and wait for export, squash, and import to finish - any errors
	// will be propagated through
	return b.loadImage(squashedImageReader, opts)
}

// error returned may be safely ignored
func (b *BuildCmd) exportImage(image *docker.APIImages, out io.Writer, opts *DockerCmdOpts) error {
	// log export starting
	b.reporter.Log(log.WithField("image_id", image.ID), "starting squash of "+image.ID[:12])

	// report docker save begin event
	b.reporter.Event(comm.EventOptions{
		EventType: comm.BuildEventSquashStartSave,
		Data: map[string]interface{}{
			"image_id": image.ID,
		},
	})

	// do docker export
	exportOpts := docker.ExportImageOptions{
		Name:         image.ID,
		OutputStream: out,
	}
	if err := opts.DockerClient.Client().ExportImage(exportOpts); err != nil {
		b.reporter.LogLevel(log.WithField("error", err), "error exporting image for squash", log.ErrorLevel)
		return err
	}

	// report event completion if successful
	b.reporter.Event(comm.EventOptions{
		EventType: comm.BuildEventSquashFinishSave,
		Data: map[string]interface{}{
			"image_id": image.ID,
		},
	})
	return nil
}

// error returned may be safely ignored
func (b *BuildCmd) squash(in *io.PipeReader, out *io.PipeWriter, retIDBuffer io.Writer) error {
	// log begin of libsquash Squash process
	b.reporter.Event(comm.EventOptions{EventType: comm.BuildEventSquashStartSquash})

	// squash image
	if err := libsquash.Squash(in, out, retIDBuffer); err != nil {

		// if error, log error
		b.reporter.LogLevel(log.WithField("error", err), "error squashing image", log.ErrorLevel)

		// if error, propagate error
		if closeerr := out.CloseWithError(err); closeerr != nil {
			b.reporter.LogLevel(log.WithField("error", closeerr), "error closing squash image write pipe", log.ErrorLevel)
			return closeerr
		}

		return err
	}

	if err := out.Close(); err != nil {
		b.reporter.LogLevel(log.WithField("error", err), "error closing squash image write pipe", log.ErrorLevel)
		return err
	}

	// report squash completed
	b.reporter.Event(comm.EventOptions{EventType: comm.BuildEventSquashFinishSquash})
	return nil
}

func (b *BuildCmd) loadImage(in *io.PipeReader, opts *DockerCmdOpts) (err error) {
	// log begin of docker load
	b.reporter.Event(comm.EventOptions{
		EventType: comm.BuildEventSquashStartLoad,
	})

	// docker load
	loadOpts := docker.LoadImageOptions{InputStream: in}
	err = opts.DockerClient.Client().LoadImage(loadOpts)

	// log completion of docker load
	b.reporter.Event(comm.EventOptions{
		EventType: comm.BuildEventSquashFinishLoad,
		Data:      map[string]interface{}{"error": err},
	})
	return
}
