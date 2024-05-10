package pipeline

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	output := in
	for _, stage := range stages {
		output = func(in In) (out Out) {
			valStream := make(Bi)
			go func() {
				defer close(valStream)
				for {
					select {
					case <-done:
						return
					case v, ok := <-in:
						if !ok {
							return
						}
						select {
						case valStream <- v:
						case <-done:
						}
					}
				}
			}()
			return stage(valStream)
		}(output)
	}
	return output
}
