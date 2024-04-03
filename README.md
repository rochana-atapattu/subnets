# Subnet Calculator

## Introduction

This Subnet Calculator is a modern CLI tool inspired by `davidc/subnets`, designed to simplify network analysis and subnet calculation tasks. Built with Go and leveraging the Bubble Tea framework for a delightful TUI (Text-based User Interface), it offers an interactive and user-friendly way to calculate subnets, understand IP address allocations, and manage network configurations efficiently.

## Features

- **Subnet Calculation:** Easily calculate subnet masks, network addresses, broadcast addresses, and available IP addresses.
- **Interactive UI:** Powered by Bubble Tea, the tool provides an interactive experience for inputting data and viewing results.
- **CIDR Support:** Full support for Classless Inter-Domain Routing (CIDR) notation to specify IP addresses and subnet masks.
- **IP Range Analysis:** Analyze and display the range of IP addresses within a given subnet.

## Installation

Ensure you have Go installed on your system. You can download and install Go from [the official Go website](https://golang.org/dl/).

To install the Subnet Calculator, run the following command:

```bash
go install github.com/rochana-atapattu/subnets@latest
```
## Usage

```bash
subnets <IP address> <mask length>
```

Example

```bash
subnets 192.168.0.0 24
```

Follow the on-screen prompts to enter your network information and perform subnet calculations.

## Contributing

We welcome contributions! If you'd like to contribute, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -am 'feat: some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by `davidc/subnets`.
- Built with the [Bubble Tea framework](https://github.com/charmbracelet/bubbletea).

