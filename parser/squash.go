package parser

import (
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/winchman/builder-core/communication"
	"github.com/winchman/libsquash"
)

// error returned may be safely ignored
func (b *BuildCmd) exportImage(image *docker.APIImages, out io.Writer, opts *DockerCmdOpts) error {
	b.reporter.Log(log.WithField("image_id", image.ID), "starting squash of "+image.ID)
	b.reporter.Event(comm.EventOptions{
		EventType: comm.BuildEventSquashStartSave,
		Data: map[string]interface{}{
			"image_id": image.ID,
		},
	})
	exportOpts := docker.ExportImageOptions{
		Name:         image.ID,
		OutputStream: out,
	}
	if err := opts.DockerClient.Client().ExportImage(exportOpts); err != nil {
		b.reporter.LogLevel(log.WithField("error", err), "error exporting image for squash", log.ErrorLevel)
		return err
	}
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
	b.reporter.Event(comm.EventOptions{EventType: comm.BuildEventSquashStartSquash})
	if err := libsquash.Squash(in, out, retIDBuffer); err != nil {
		b.reporter.LogLevel(log.WithField("error", err), "error squashing image", log.ErrorLevel)
		if err := out.CloseWithError(err); err != nil {
			b.reporter.LogLevel(log.WithField("error", err), "error closing squash image write pipe", log.ErrorLevel)
			return err
		}
		return err
	} else {
		if err := out.Close(); err != nil {
			b.reporter.LogLevel(log.WithField("error", err), "error closing squash image write pipe", log.ErrorLevel)
			return err
		}
	}
	b.reporter.Event(comm.EventOptions{EventType: comm.BuildEventSquashFinishSquash})
	return nil
}

func (b *BuildCmd) loadImage(in *io.PipeReader, opts *DockerCmdOpts) error {
	b.reporter.Event(comm.EventOptions{EventType: comm.BuildEventSquashStartLoad})
	loadOpts := docker.LoadImageOptions{InputStream: in}
	return opts.DockerClient.Client().LoadImage(loadOpts)
}
