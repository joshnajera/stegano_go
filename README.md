#Created by: Joshua Najera

About
----
This is a simple program that hides data in the least significant bit (LSB) of each color for each pixel. The **bottom right 11 pixels** are used to store the message length. The **remaining pixels** are used for the message itself.

####How this works:

Starting from the last pixel, read the LSB of the **RED** value, *then* that of the **GREEN** value, and *finally* that of the **BLUE** value. When finished, move to the pixel in front of the current one and repeat. Repeat this process as needed. When put together this forms the binary representation of the message and its length.


Requirements
----
You must have [golang](https://golang.org/) installed to run/build from source code 

Usage
----
To run golang code from source you can use the following from terminal:

    go run main.go

Alternatively you can build using

    go build main.go

And then run using the created binary

#####Example

    go run main.go -r image.png

or

    main.exe -r image.png

Flags & Arguments
----
Writing a message**

    -w imageName "message here"

Writing from text file

    -f imageeName "file name"

Reading from an image

    -r imageeName