http://brycekerley.net/blog/2009/06/trivia300.html

I was playing with some code to solve the problem described above, and...

    git checkout working
    6g mazes.go && 6l -o mazes mazes.6
    ./mazes mazes.txt

Yields:

    <snip 365 width/walker messages>
    throw: all goroutines are asleep - deadlock!

This is the expected behaviour, as I've not put error detection code in for the
problem maze 366 exhibits -- s and f are not connected by any path.

However, enabling runtime.GOMAXPROCS(2) causes some errors as follows:

    git checkout broken
    6g mazes.go && 6l -o mazes mazes.6
    ./mazes mazes.txt

Yields (note normal output is disabled as it tramples over the runtime error):

    throw: assert
    panic PC=0x7facf21d5d98
    throw+0x3e /home/alex/go/src/pkg/runtime/runtime.c:74
        throw(0x455ad1, 0x0)
    hash_insert_internal+0x3fa /home/alex/go/src/pkg/runtime/hashmap.c:416
        hash_insert_internal(0xf2190258, 0x7fac, 0x2, 0x7fac, 0x2a48700, ...)
    hash_grow+0xea /home/alex/go/src/pkg/runtime/hashmap.c:283
        hash_grow(0xf2190230, 0x7fac, 0xf2190258, 0x7fac, 0x2, ...)
    hash_insert_internal+0x19c /home/alex/go/src/pkg/runtime/hashmap.c:450
        hash_insert_internal(0xf2190258, 0x7fac, 0x0, 0x7fac, 0xa4cf75c0, ...)
    hash_insert+0x6e /home/alex/go/src/pkg/runtime/hashmap.c:460
        hash_insert(0xf2190230, 0x7fac, 0xf21d6050, 0x7fac, 0xf21d5ff8, ...)
    mapassign+0x5e /home/alex/go/src/pkg/runtime/hashmap.c:821
        mapassign(0xf2190230, 0x7fac, 0xf21d6050, 0x7fac, 0xf21d6058, ...)
    runtime·mapassign1+0x4b /home/alex/go/src/pkg/runtime/hashmap.c:849
        runtime·mapassign1(0xf2190230, 0x7fac)
    main·*walker·walk+0x8e /home/alex/git/mazes/mazes.go:197
        main·*walker·walk(0xf218ab40, 0x7fac, 0xf218e240, 0x7fac, 0xf218e300, ...)
    goexit /home/alex/go/src/pkg/runtime/proc.c:136
        goexit()
    0x7facf218ab40 unknown pc

    goroutine 3 [4]:
    gosched+0x34 /home/alex/go/src/pkg/runtime/proc.c:522
        gosched()
    chanrecv+0x167 /home/alex/go/src/pkg/runtime/chan.c:347
        chanrecv(0xf218d1e0, 0x7fac, 0xf21951a8, 0x7fac, 0x0, ...)
    runtime·chanrecv1+0x50 /home/alex/go/src/pkg/runtime/chan.c:417
        runtime·chanrecv1(0xf218d1e0, 0x7fac)
    main·readMazes+0xe0 /home/alex/git/mazes/mazes.go:112
        main·readMazes(0xf218d1e0, 0x7fac, 0xf218d240, 0x7fac)
    goexit /home/alex/go/src/pkg/runtime/proc.c:136
        goexit()
    0x7facf218d1e0 unknown pc

    goroutine 2 [4]:
    gosched+0x34 /home/alex/go/src/pkg/runtime/proc.c:522
        gosched()
    chansend+0x14c /home/alex/go/src/pkg/runtime/chan.c:245
        chansend(0xf218d1e0, 0x7fac, 0xf2194070, 0x7fac, 0x0, ...)
    runtime·chansend1+0x54 /home/alex/go/src/pkg/runtime/chan.c:389
        runtime·chansend1(0xf218d1e0, 0x7fac)
    main·_func_001+0x122 /home/alex/git/mazes/mazes.go:85
        main·_func_001(0xf2189018, 0x7fac, 0xf218c790, 0x7fac, 0xf2189020, ...)
    goexit /home/alex/go/src/pkg/runtime/proc.c:136
        goexit()

    goroutine 1 [4]:
    gosched+0x34 /home/alex/go/src/pkg/runtime/proc.c:522
        gosched()
    chanrecv+0x167 /home/alex/go/src/pkg/runtime/chan.c:347
        chanrecv(0xf218d240, 0x7fac, 0xf2192f40, 0x7fac, 0x0, ...)
    runtime·chanrecv1+0x50 /home/alex/go/src/pkg/runtime/chan.c:417
        runtime·chanrecv1(0xf218d240, 0x7fac)
    main·main+0x126 /home/alex/git/mazes/mazes.go:49
        main·main()
    mainstart+0xf /home/alex/go/src/pkg/runtime/amd64/asm.s:54
        mainstart()
    goexit /home/alex/go/src/pkg/runtime/proc.c:136
        goexit()
