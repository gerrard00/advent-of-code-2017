# Advent Of Code 2017

I'm using Advent of Code 2017 as an excuse to learn some golang. I put this up so I can work on it on any computer. These are not the droids you are looking for. For which you are looking?

## Stuff I did wrong

* Figured out the spiral thing with math so I didn't have to generate the spiral, then I still had to generate the spiral for part two.

* wrote code that used a stack to do `stack = stack.pop()` because I couldn't figure out how to use pointers to pointers `**` with a struct method. I could have made this work without having to constantly do assignments by just using normal function that accepted pointers to pointers (i.e. `func foo(**stack)`).
