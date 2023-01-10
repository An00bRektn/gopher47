# gopher47
> A 3rd party implant for the [Havoc](https://github.com/HavocFramework/Havoc) Framework

<p align="center">
    <img src="assets/gopher47.png">
</p>

<sub align="right">it's like the videogame :D</sub>

## What is This? 🤔
This is a 3rd Party agent for the Havoc C2 written in Golang, mainly written as a learning project, but I'm sure it's still useful.

By the time this is out, you can read more about what and why this is at my blog: [here](https://notateamserver.xyz)

**Contributions welcome!** I don't plan on updating this all too regularly, but if I have fun making it I might add more stuff over time.

## Commands

| Command | Description                                     | Example            |
|---------|-------------------------------------------------|--------------------|
| `o7`    | The gopher dies :(                              | o7                 |
| `shell` | Run a command (executed through Go's `os/exec`) | shell netstat -ano |
| `kill`  | Kills a process by PID                          | kill 31337         |


## Usage
Once you have your teamserver up, it's as simple as running the following:
- **Attacking Machine**: `python handler.py`
- **Target Machine**: `./gopher47`

You can use the Havoc GUI to compile it, or you can just edit the source code as you please and play with the Makefile, there isn't that much of it.

## FAQ

### Why Go?
I just wanted to have an actual Golang project put together that I can [point to](https://i.kym-cdn.com/entries/icons/original/000/035/627/cover2.jpg).

### Will it evade AV/EDR?
idk, but grow up. Obfuscate and customize it yourself, stop being a baby.

### Why's the binary so large though?
Golang, along with Rust and other languages, compile **statically**, meaning all of the libraries necessary to run the executable are baked into the binary, which adds up. If you want to reduce the size, I won't do it by default, but check out [this link](https://github.com/xaionaro/documentation/blob/master/golang/reduce-binary-size.md) for some tips. 

### How's your day going?
I photoshopped a gun into the Golang gopher's hand for this at 1:00 AM, and my Winter break is over tomorrow.


## Acknowledgements/References/Related Work
- [C5pider](https://github.com/Cracked5pider) and the entire [Havoc Framework](https://github.com/HavocFramework) team for letting me not have to write my own C2 to do some implant development
- [SharpAgent](https://github.com/susMdT/SharpAgent/) and [PyHmmm](https://github.com/CodeXTF2/PyHmmm) were great ~~projects to steal from~~ reference material
- [OffensiveGolang](https://github.com/bluesentinelsec/OffensiveGoLang) has some neat Offensive Go work
- [maldev-for-dummies](https://github.com/chvancooten/maldev-for-dummies) also is a nice starting point for working with not C malware

