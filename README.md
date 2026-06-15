# qoin-shor

## Project Introduction
**qoin-shor** is a Decentralized Quantum-resistant AI Supercomputing Network.

## Building & Documentation

## Contributing

### Basic Build Instructions

#### System-specific Software Dependencies
Building Shor requires some system dependencies, usually provided by your distribution.

**Ubuntu/Debian:**
```bash
sudo apt install mesa-opencl-icd ocl-icd-opencl-dev gcc git bzr jq pkg-config curl clang build-essential hwloc libhwloc-dev wget -y && sudo apt upgrade -y
```

**Fedora:**
```bash
sudo dnf -y install gcc make git bzr jq pkgconfig mesa-libOpenCL mesa-libOpenCL-devel opencl-headers ocl-icd ocl-icd-devel clang llvm wget hwloc hwloc-devel
```

For other distributions you can find the required dependencies here. For instructions specific to macOS, you can find them here.

#### Go
To build Shor, you need a working installation of Go 1.22.2 or higher:
```bash
wget -c https://golang.org/dl/go1.22.2.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
```

**TIP:** You'll need to add `/usr/local/go/bin` to your path. For most Linux distributions you can run something like:
```bash
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc && source ~/.bashrc
```
See the official Golang installation instructions if you get stuck.

## Build and install Shor

1. Once all the dependencies are installed, you can build and install the Shor suite.
   
   Clone the repository:
   ```bash
   git clone https://github.com/qoin/qoin-shor.git
   cd qoin-shor
   ```

2. To join mainnet, checkout the latest release.
   ```bash
   git checkout <latest_release_tag>
   ```

3. Build and install Shor:
   ```bash
   make clean all
   sudo make install
   ```

4. You should now have Shor installed. You can now start the Shor daemon and sync the chain. 
   
   **Command Example:**
   ```bash
   shor daemon
   ```

## License

Dual-licensed under MIT + Apache 2.0.
