# Slocc 🦥

## Description
Slocc is a clone of David Wheeler's [sloccount](https://dwheeler.com/sloccount/) written in Go. The idea is to provide a similar interface and output while offering easy configuration and extensibility.

## Why this
Well, I like `sloccount`. It provides some sort of vanity metrics when I'm working on my projects, but it also offers me some quick insights when I'm examining new codebases. `sloccount` is fine, so why make yet another one? The answer is simple: "Why not?". I guess my main pain points with `sloccount` were that it didn't support off-the-batch the languages I was using, I could configure it easily, and I had to do that every time I spun up a fresh GNU/Linux on a new machine. The other main reason is that, despite being somewhat familiar with COCOMO, it didn't matter much to me most of the time, so I wanted to hack the output. 
Yes, I considered forking `sloccount` and creating my own version, but then I thought, why not redo it in Go instead? ... So, here it is. Still a work in progress, but the main functionality is out there.

## Quick start
Clone the project and use the Makefile as follows to build the project:
```
make 
```
It should compile and generate a binary: `slocc`. 

## Usage
Then, to use it, simply type `slocc` followed by the directory or file where your source code is located.
```
slocc my_project_dir
```

## Contributing
I don't have big plans for the project. However, I'll likely continue to add a few features every now and then.
It would be fun to hear you want to add or change something, so go ahead, create a PR, or DM me :-)
