package shared

type NautilusRPCServer struct {
	Impl NautilusPipelineProcessor
}

func (s *NautilusRPCServer) Name(_ struct{}, resp *string) error {
	result, err := s.Impl.Name()
	if err != nil {
		return err
	}

	*resp = result
	return nil
}

func (s *NautilusRPCServer) Process(args *NautilusPipelineContext, resp *NautilusPipelineContext) error {
	result, err := s.Impl.Process(args)
	if err != nil {
		return err
	}

	*resp = *result
	return nil
}

func (s *NautilusRPCServer) Configure(_ struct{}, resp *error) error {
	return s.Impl.Configure()
}

func (s *NautilusRPCServer) Health(_ struct{}, resp *error) error {
	return s.Impl.Health()
}
