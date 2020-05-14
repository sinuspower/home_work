package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	I   = interface{}
	In  = <-chan I
	Out = In
	Bi  = chan I
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
put := func(in In, out Bi) { // puts value from in channel to out or stops working by done channel
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok { // input channel closed
					return
				}
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}

	valueStream := in
	for _, stage := range stages {
		stageIn := make(Bi)
		go put(valueStream, stageIn)
		valueStream = stage(stageIn)
	}

	return valueStream
}
