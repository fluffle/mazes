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
    panic PC=0x7f4c68b43030
    throw+0x3e /home/alex/go/src/pkg/runtime/runtime.c:74
        throw(0x434a12, 0x0)
    hash_grow+0xfd /home/alex/go/src/pkg/runtime/hashmap.c:284
        hash_grow(0x68bcfcb0, 0x7f4c, 0x68bcfcd8, 0x7f4c, 0x2, ...)
    hash_insert_internal+0x19c /home/alex/go/src/pkg/runtime/hashmap.c:450
        hash_insert_internal(0x68bcfcd8, 0x7f4c, 0x0, 0x7f4c, 0x45366480, ...)
    hash_insert+0x6e /home/alex/go/src/pkg/runtime/hashmap.c:460
        hash_insert(0x68bcfcb0, 0x7f4c, 0x68b43250, 0x7f4c, 0x68b431f8, ...)
    mapassign+0x5e /home/alex/go/src/pkg/runtime/hashmap.c:821
        mapassign(0x68bcfcb0, 0x7f4c, 0x68b43250, 0x7f4c, 0x68b43258, ...)
    runtime·mapassign1+0x4b /home/alex/go/src/pkg/runtime/hashmap.c:849
        runtime·mapassign1(0x68bcfcb0, 0x7f4c)
    main·*walker·walk+0x8e /home/alex/git/mazes/mazes.go:195
        main·*walker·walk(0x68b789a0, 0x7f4c, 0x68bb31b0, 0x7f4c, 0x68bb3570, ...)
    goexit /home/alex/go/src/pkg/runtime/proc.c:136
        goexit()
    0x7f4c68b789a0 unknown pc

    goroutine 3 [2]:
    gosched+0x34 /home/alex/go/src/pkg/runtime/proc.c:522
        gosched()
    chanrecv+0x179 /home/alex/go/src/pkg/runtime/chan.c:349
        chanrecv(0x68b3a1e0, 0x7f4c, 0x68b421a8, 0x7f4c, 0x0, ...)
    runtime·chanrecv1+0x50 /home/alex/go/src/pkg/runtime/chan.c:417
        runtime·chanrecv1(0x68b3a1e0, 0x7f4c)
    main·readMazes+0xe0 /home/alex/git/mazes/mazes.go:110
        main·readMazes(0x68b3a1e0, 0x7f4c, 0x68b3a240, 0x7f4c)
    goexit /home/alex/go/src/pkg/runtime/proc.c:136
        goexit()
    0x7f4c68b3a1e0 unknown pc

    goroutine 2 [1]:
    gosched+0x34 /home/alex/go/src/pkg/runtime/proc.c:522
        gosched()
    chansend+0x14c /home/alex/go/src/pkg/runtime/chan.c:245
        chansend(0x68b3a1e0, 0x7f4c, 0x68b41068, 0x7f4c, 0x0, ...)
    runtime·chansend1+0x54 /home/alex/go/src/pkg/runtime/chan.c:389
        runtime·chansend1(0x68b3a1e0, 0x7f4c)
    main·_func_001+0xfa /home/alex/git/mazes/mazes.go:84
        main·_func_001(0x68b36028, 0x7f4c, 0x68b36020, 0x7f4c, 0x68b36018, ...)
    goexit /home/alex/go/src/pkg/runtime/proc.c:136
        goexit()

    goroutine 1 [4]:
    gosched+0x34 /home/alex/go/src/pkg/runtime/proc.c:522
        gosched()
    chanrecv+0x167 /home/alex/go/src/pkg/runtime/chan.c:347
        chanrecv(0x68b3a240, 0x7f4c, 0x68b3ff40, 0x7f4c, 0x0, ...)
    runtime·chanrecv1+0x50 /home/alex/go/src/pkg/runtime/chan.c:417
        runtime·chanrecv1(0x68b3a240, 0x7f4c)
    main·main+0x126 /home/alex/git/mazes/mazes.go:49
        main·main()
    mainstart+0xf /home/alex/go/src/pkg/runtime/amd64/asm.s:54
        mainstart()
    goexit /home/alex/go/src/pkg/runtime/proc.c:136
        goexit()

I did also see once:

    throw: throw: assert
    panic PC=0x7f302fe56e30
    assert
    double panic

Which was more worrying! OH NOES DOUBLE PANIC!

