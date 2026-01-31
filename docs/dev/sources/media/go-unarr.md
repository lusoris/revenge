# gen2brain/go-unarr

> Source: https://pkg.go.dev/github.com/gen2brain/go-unarr
> Fetched: 2026-01-30T23:54:18.670728+00:00
> Content-Hash: 29c8d19a0896602b
> Type: html

---

Overview

¶

Package unarr is a decompression library for RAR, TAR, ZIP and 7z archives.

Index

¶

Variables

type Archive

func NewArchive(path string) (a *Archive, err error)

func NewArchiveFromMemory(b []byte) (a *Archive, err error)

func NewArchiveFromReader(r io.Reader) (a *Archive, err error)

func (a *Archive) Close() (err error)

func (a *Archive) Entry() error

func (a *Archive) EntryAt(off int64) error

func (a *Archive) EntryFor(name string) error

func (a *Archive) Extract(path string) (contents []string, err error)

func (a *Archive) List() (contents []string, err error)

func (a *Archive) ModTime() time.Time

func (a *Archive) Name() string

func (a *Archive) Offset() int64

func (a *Archive) RawName() string

func (a *Archive) Read(b []byte) (n int, err error)

func (a *Archive) ReadAll() ([]byte, error)

func (a *Archive) Seek(offset int64, whence int) (int64, error)

func (a *Archive) Size() int

Constants

¶

This section is empty.

Variables

¶

View Source

var (

ErrOpenFile    =

errors

.

New

("unarr: open file failed")

ErrOpenMemory  =

errors

.

New

("unarr: open memory failed")

ErrOpenArchive =

errors

.

New

("unarr: no valid RAR, ZIP, 7Z or TAR archive")

ErrEntry       =

errors

.

New

("unarr: failed to parse entry")

ErrEntryAt     =

errors

.

New

("unarr: failed to parse entry at")

ErrEntryFor    =

errors

.

New

("unarr: failed to parse entry for")

ErrSeek        =

errors

.

New

("unarr: seek failed")

ErrRead        =

errors

.

New

("unarr: read failure")

)

Functions

¶

This section is empty.

Types

¶

type

Archive

¶

type Archive struct {

// contains filtered or unexported fields

}

Archive represents unarr archive

func

NewArchive

¶

func NewArchive(path

string

) (a *

Archive

, err

error

)

NewArchive returns new unarr Archive

func

NewArchiveFromMemory

¶

func NewArchiveFromMemory(b []

byte

) (a *

Archive

, err

error

)

NewArchiveFromMemory returns new unarr Archive from byte slice

func

NewArchiveFromReader

¶

func NewArchiveFromReader(r

io

.

Reader

) (a *

Archive

, err

error

)

NewArchiveFromReader returns new unarr Archive from io.Reader

func (*Archive)

Close

¶

func (a *

Archive

) Close() (err

error

)

Close closes the underlying unarr archive

func (*Archive)

Entry

¶

func (a *

Archive

) Entry()

error

Entry reads the next archive entry.

io.EOF is returned when there is no more to be read from the archive.

func (*Archive)

EntryAt

¶

func (a *

Archive

) EntryAt(off

int64

)

error

EntryAt reads the archive entry at the given offset

func (*Archive)

EntryFor

¶

func (a *

Archive

) EntryFor(name

string

)

error

EntryFor reads the (first) archive entry associated with the given name

func (*Archive)

Extract

¶

func (a *

Archive

) Extract(path

string

) (contents []

string

, err

error

)

Extract extracts archive to destination path

func (*Archive)

List

¶

func (a *

Archive

) List() (contents []

string

, err

error

)

List lists the contents of archive

func (*Archive)

ModTime

¶

func (a *

Archive

) ModTime()

time

.

Time

ModTime returns the stored modification time of the current entry

func (*Archive)

Name

¶

func (a *

Archive

) Name()

string

Name returns the name of the current entry as UTF-8 string

func (*Archive)

Offset

¶

func (a *

Archive

) Offset()

int64

Offset returns the stream offset of the current entry, for use with EntryAt

func (*Archive)

RawName

¶

added in

v0.1.2

func (a *

Archive

) RawName()

string

RawName returns the name of the current entry as raw string

func (*Archive)

Read

¶

func (a *

Archive

) Read(b []

byte

) (n

int

, err

error

)

Read tries to read 'b' bytes into buffer, advancing the read offset pointer.

Returns the actual number of bytes read.

func (*Archive)

ReadAll

¶

func (a *

Archive

) ReadAll() ([]

byte

,

error

)

ReadAll reads current entry and returns data

func (*Archive)

Seek

¶

func (a *

Archive

) Seek(offset

int64

, whence

int

) (

int64

,

error

)

Seek moves the read offset pointer interpreted according to whence.

Returns the new offset.

func (*Archive)

Size

¶

func (a *

Archive

) Size()

int

Size returns the total size of uncompressed data of the current entry