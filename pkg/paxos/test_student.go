package paxos

import (
	"coms4113/hw5/pkg/base"
)

// Fill in the function to lead the program to a state where A2 rejects the Accept Request of P1
func ToA2RejectP1() []func(s *base.State) bool {

	p3PreparePhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		check := s3.proposer.Phase == Propose
		if check {
			//fmt.Printf("s3 enters Propose phase with N: %d\n", s3.proposer.N)
		}
		return check
	}
	p1PreparePhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		check := s1.proposer.Phase == Propose
		if check {
			//fmt.Printf("s1 enters Propose phase with N: %d\n", s1.proposer.N)
		}
		return check
	}

	p1AcceptPhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		check := s1.proposer.Phase == Accept
		if check {
			//fmt.Printf("s1 enters Accept phase with n_p: %d, v_a: %+v\n", s1.n_p, s1.v_a)
		}
		return check
	}

	a2RejectAccept := func(s *base.State) bool {
		s2 := s.Nodes()["s2"].(*Server)
		s1 := s.Nodes()["s1"].(*Server)
		check := s1.proposer.Phase == Accept && s1.proposer.N < s2.n_p
		if check {
			//fmt.Printf("s2 rejects enters Accept phase with n_p: %d, v_a: %+v\n", s2.n_p, s2.v_a)
		}
		return check
	}

	return []func(s *base.State) bool{p1PreparePhase, p3PreparePhase, p1AcceptPhase, a2RejectAccept}
	//panic("fill me in")
}

// Fill in the function to lead the program to a state where a consensus is reached in Server 3.
func ToConsensusCase5() []func(s *base.State) bool {

	p3AcceptPhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		check := s3.proposer.Phase == Accept
		if check {
			//fmt.Printf("s3 enters Accept phase with n_p: %d, v_a: %+v\n", s3.n_p, s3.v_a)
		}
		return check
	}

	s3KnowConsensus := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		check := s3.agreedValue == "v3"
		return check
	}

	return []func(s *base.State) bool{p3AcceptPhase, s3KnowConsensus}

	//panic("fill me in")
}

// Fill in the function to lead the program to a state where all the Accept Requests of P1 are rejected
func NotTerminate1() []func(s *base.State) bool {

	p1PreparePhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		check := s1.proposer.Phase == Propose
		if check {
			//fmt.Printf("s1 enters Propose phase with N: %d\n", s1.proposer.N)
		}
		return check
	}

	p3PreparePhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		s1 := s.Nodes()["s1"].(*Server)
		check := s3.proposer.Phase == Propose && s3.proposer.N == s1.proposer.N+1 && s3.proposer.N == s1.n_p && s1.proposer.Phase == Propose
		if check {
			//fmt.Printf("s3 enters Propose phase with N: %d, ResponseCount: %d\n", s3.proposer.N, s3.proposer.ResponseCount)
		}
		return check
	}

	// Accept Phase
	p1AcceptPhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		s3 := s.Nodes()["s3"].(*Server)
		check := s1.proposer.Phase == Accept && s3.proposer.Phase == Propose //&& s1.proposer.N < s3.n_p && s1.proposer.SuccessCount == 0 && s1.proposer.ResponseCount == 3
		if check {
			//fmt.Printf("s1 AcceptFailed phase with N: %d, ResponseCount: %d\n", s1.proposer.N, s1.proposer.ResponseCount)
		}
		return check
	}

	a2RejectAccept := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		s2 := s.Nodes()["s2"].(*Server)
		s1 := s.Nodes()["s1"].(*Server)
		check := s1.proposer.Phase == Accept && s3.proposer.Phase == Propose && s1.proposer.N < s3.n_p && s1.proposer.N < s2.n_p && s1.proposer.N < s1.n_p
		if check {
			//fmt.Printf("s1 rejects  with n_p: %d, v_a: %+v\n", s3.n_p, s3.v_a)
		}
		return check
	}

	p1RejectAccept := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		check := s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 1
		if check {
			//fmt.Printf("s1 reject with n_p: %d, v_a: %+v\n", s1.n_p, s1.v_a)
		}
		return check
	}

	return []func(s *base.State) bool{p1PreparePhase, p3PreparePhase, p1AcceptPhase, a2RejectAccept, p1RejectAccept}
	//panic("fill me in")
}

// Fill in the function to lead the program to a state where all the Accept Requests of P3 are rejected
func NotTerminate2() []func(s *base.State) bool {

	p3PreparePhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		s1 := s.Nodes()["s1"].(*Server)
		check := s1.proposer.Phase == Propose && s1.proposer.N == s3.proposer.N+1 && s3.proposer.Phase == Accept
		if check {
			//fmt.Printf("s1 enters Propose phase with N: %d, ResponseCount: %d\n", s1.proposer.N, s1.proposer.ResponseCount)
		}
		return check
	}

	a2RejectAccept := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		s2 := s.Nodes()["s2"].(*Server)
		s1 := s.Nodes()["s1"].(*Server)
		check := s3.proposer.Phase == Accept && s1.proposer.Phase == Accept && s3.proposer.N < s3.n_p && s3.proposer.N < s2.n_p && s3.proposer.N < s1.n_p
		if check {
			//fmt.Printf("s3 rejects accepts N: %d, V: %+v, ResponseCount: %d, s1.n_p: %d\n", s3.proposer.N, s3.proposer.V, s3.proposer.ResponseCount, s1.n_p)
		}
		return check
	}

	p1RejectAccept := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		check := s3.proposer.Phase == Accept && s3.proposer.ResponseCount == 1 && s3.proposer.SuccessCount == 0
		if check {
			//fmt.Printf("s3 rejects accepts N: %d, V: %+v, ResponseCount: %d, SuccessCount: %d, s3.n_p: %d\n", s3.proposer.N, s3.proposer.V, s3.proposer.ResponseCount, s3.proposer.SuccessCount, s3.n_p)
		}
		return check
	}

	return []func(s *base.State) bool{p3PreparePhase, a2RejectAccept, p1RejectAccept}
	//panic("fill me in")
}

// Fill in the function to lead the program to a state where all the Accept Requests of P1 are rejected again.
func NotTerminate3() []func(s *base.State) bool {

	return NotTerminate1()
	//panic("fill me in")
}

// Fill in the function to lead the program to make P1 propose first, then P3 proposes, but P1 get rejects in
// Accept phase
func concurrentProposer1() []func(s *base.State) bool {

	p1PreparePhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		check := s1.proposer.Phase == Propose
		if check {
			//fmt.Printf("s1 enters Propose phase with N: %d\n", s1.proposer.N)
		}
		return check
	}

	p3PreparePhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		s1 := s.Nodes()["s1"].(*Server)
		check := s3.proposer.Phase == Propose && s3.proposer.N == s1.proposer.N+1
		if check {
			//fmt.Printf("s1 enters Propose phase with N: %d\n", s3.proposer.N)
		}
		return check
	}

	p1AcceptPhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		s3 := s.Nodes()["s3"].(*Server)
		check := s1.proposer.Phase == Accept && s3.proposer.Phase == Propose
		if check {
			//fmt.Printf("s1 enters Accept phase with n_p: %d, v_a: %+v\n", s1.n_p, s1.v_a)
		}
		return check
	}

	p3AcceptPhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		s1 := s.Nodes()["s1"].(*Server)
		s2 := s.Nodes()["s2"].(*Server)
		check := s3.proposer.Phase == Accept && s1.proposer.N+1 == s3.proposer.N && s1.proposer.N < s1.n_p && s1.proposer.N < s2.n_p && s1.proposer.N < s3.n_p
		if check {
			//fmt.Printf("s3 enters Accept phase with n_p: %d, v_a: %+v\n", s3.n_p, s3.v_a)
		}
		return check
	}

	p1RejectFailed := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		check := s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 1 && s1.proposer.SuccessCount == 0
		if check {
			//fmt.Printf("s3 enters Accept phase with n_p: %d, v_a: %+v\n", s1.n_p, s1.v_a)
		}
		return check
	}

	return []func(s *base.State) bool{p1PreparePhase, p3PreparePhase, p1AcceptPhase, p3AcceptPhase, p1RejectFailed}

	//panic("fill me in")
}

// Fill in the function to lead the program continue  P3's proposal  and reaches consensus at the value of "v3".
func concurrentProposer2() []func(s *base.State) bool {

	p3AcceptPhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		s2 := s.Nodes()["s2"].(*Server)
		s1 := s.Nodes()["s1"].(*Server)
		check := s3.proposer.Phase == Accept && s1.n_p == s3.proposer.N && s2.n_p == s3.proposer.N && s3.n_p == s3.proposer.N && s3.proposer.V == "v3"
		if check {
			//fmt.Printf("s1 enters Propose phase with N: %d, ResponseCount: %d\n", s1.proposer.N, s1.proposer.ResponseCount)
		}
		return check
	}

	p3DecidePhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		check := s3.proposer.Phase == Decide && s3.agreedValue == "v3"
		if check {
			//fmt.Printf("s1 enters Propose phase with N: %d, ResponseCount: %d\n", s3.proposer.N, s3.proposer.ResponseCount)
		}
		return check
	}

	return []func(s *base.State) bool{p3AcceptPhase, p3DecidePhase}
	//panic("fill me in")
}
