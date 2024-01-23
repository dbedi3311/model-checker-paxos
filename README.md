# Paxos Model Checker - Distributed Systems Verification Project

## Overview

The Paxos Model Checker project represents a comprehensive effort to implement and verify the Paxos algorithm, a fundamental protocol in distributed systems for achieving consensus. This project is developed in Go and focuses on model checking, a method used to systematically explore all possible states of a distributed system to ensure correctness and reliability, even in rare and challenging scenarios.

### Project Description

- **Model Checking in Distributed Systems**: Implemented a model checker from scratch to simulate and test distributed systems. This approach goes beyond traditional fuzz testing by exploring all possible system states.
- **Paxos Protocol Implementation**: Developed an implementation of the Paxos consensus protocol, a core algorithm in distributed computing. Ensured adherence to the protocol's specifications and behaviors.

### Features

- **State Machine Abstraction**: The entire system (nodes + network) is abstracted as a state machine. This includes handling events like message arrivals and timer triggers. Provides an interface for specifying system's state machine, defining test scenarios, and running BFS algorithm to explore all possible events and interleavings up to final states while checking invariants.
- **Event-Driven Architecture**: The system evolves by processing events such as message delivery and timer activations, leading to new system states.
- **Comprehensive Testing**: Utilized Breadth First Search (BFS) to walk through all states, validating the correctness of the Paxos implementation under various scenarios. 

### Challenges and Solutions

- **Handling Network Abstractions**: Modeling network behaviors, such as message drops and delays, was challenging. A detailed abstraction of on-the-fly messages was implemented to simulate real network conditions.
- **State Explosion**: The challenge of state explosion was addressed by carefully designing test scenarios and using efficient state representation and search strategies.
- **Ensuring Immutability**: Maintaining immutability in nodes and messages was crucial for accurate state representation. This was achieved by creating copies of nodes and messages before modifications.

### Paxos Implementation

- **Node and Timer Interfaces**: Implemented interfaces for nodes and timers, forming the basis of the state machine.
- **Message Handling**: Developed a sophisticated message handling system to simulate various network conditions and node behaviors.
- **Proposer and Acceptor Logic**: Carefully implemented the logic for proposers and acceptors in the Paxos protocol, ensuring correct state transitions and consensus achievement.

### Conclusion

This project showcases a practical implementation of model checking in verifying the Paxos protocol, addressing complex challenges inherent in distributed systems. It serves as a valuable reference for understanding and verifying distributed consensus algorithms.

## Technologies

- **Language**: Go (version 1.19.1)
- **Testing Frameworks**: Utilized Go's built-in testing frameworks for comprehensive unit and integration tests.

## How to Run




