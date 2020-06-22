# BytePointerEncode: Golang Edition

The unexpected sequel for the project i've left in dust.

## What is this

It's in the title bro.

## What's the difference

This one is rewritten in Golang.
I heard Golang was like a lovechild between C and Python (in terms of performance and readability at least), and hoo boi it is.
The only thing i kinda don't like is that it's clunky, but not as clunky as C so i'm grateful for that.

Also i have 3 version of it.
The basic one is like the standard way to use it. Not Unoptimized as hell, but easily understandable (i hope).
The compression one is, well, the one which will compress the output file. Cause mind you, the output can be as big as 8x the original size. This one will automatically compress the output if you have no spare space. But expect the memory usage to go up. 
The bufio one is incomplete. Somehow the bufio.read function didn't fill the buffer in some case, and messes up my program. Still trying to workaround this one.

## Todo

- Make a friendly-to-non-nerd version

## Can i contribute

Sure you can bro
