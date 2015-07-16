package node

import "github.com/cloud66/starter/packs"

type DockerfileContext struct {
	packs.DockerfileContextBase
}

type DockerfileWriter struct {
	packs.DockerfileWriterBase
}

func (w *DockerfileWriter) Write(context *DockerfileContext) error {
	return w.DockerfileWriterBase.Write(context)
}
