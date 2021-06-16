# Yaegi Interpreter

Yaegi is a Go interpreter. It powers executable Go scripts and plugins, in embedded interpreters or interactive shells, on top of Go runtime.

# Usage

Yaegi github readme has a nice write up about the usage of yaegi: [Link](https://github.com/traefik/yaegi#as-a-dynamic-extension-framework)

For exposing an interface to the plugin i.e the code we're interpreting, it is essential to generate a symbol table. Interfaces cannot be added dynamically, it is required to pre-compile the interface wrappers. 

The interface wrappers can be generated using the `yaegi extract` command. The wrapper generated from the `yaegi extract` can be used via the `interp.Use` api/method. Hence for each interface exposed as a plugin it is necessary to generate interface wrappers. 

# Pros

1. Speed: Fastest among all the plugin systems outside of go native plugin.
2. No hassle of handling plugin binaries. The plugin code is interpreted directly on the fly, hence there isn't any need to build the binaries.
3. Simple interpreter API: `New()`, `Eval()`, `Use()`.

# Cons 

1. Does not support Go modules
2. Bug: Panics on importing Protobuf.
3. Under active development (Stable version not yet released)
4. Not enough documentation around it.