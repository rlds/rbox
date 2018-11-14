package bworker

type Bworker interface {
	Run(in map[string]string) (rets string, datatype string)
}
