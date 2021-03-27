# note

A simple CLI note taking app. Still WIP as I keep extending it to meet my
needs.

Well, technically the note taking is done with your favorite editor (`export
EDITOR=vim`) - not this app. This app is mostly about note discovery and is
inspired by the Zettelkasten (https://en.wikipedia.org/wiki/Zettelkasten) note
taking method.

How does this work? All notes are plain-text files on your system. You simply
create new notes with `note new`. I like writing them as plain-text files or in
markdown. Each note will get a unique identifier. The notes can be linked
together with `[[<id>]]` e.g. `[[20210324085947]]`. These links are later on
displayed in various forms during note search process.

It's recommended to create new notes frequently, keep them short and to the
point. Then use `note check` tool to make sure that notes are linked together.
This let you slwoly build and maintain a knowledge base. At least that's the
intention.

## installation

You will need Go compiler to build this. Then just run `make` which should
produce the binary.
