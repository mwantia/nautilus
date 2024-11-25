package shared

type RpcServer struct {
	Impl PipelineProcessor
}

func (s *RpcServer) Name(_ struct{}, resp *string) error {
	result, err := s.Impl.Name()
	if err != nil {
		return err
	}

	*resp = result
	return nil
}

func (s *RpcServer) Process(args *PipelineContextData, resp *PipelineContextData) error {
	result, err := s.Impl.Process(args)
	if err != nil {
		return err
	}

	*resp = *result
	return nil
}

func (s *RpcServer) Configure(cfg map[string]interface{}, resp *error) error {
	return s.Impl.Configure(cfg)
}

func (s *RpcServer) Health(_ struct{}, resp *error) error {
	return s.Impl.Health()
}
