# gen2brain/go-unarr

> Source: https://pkg.go.dev/github.com/gen2brain/go-unarr
> Fetched: 2026-01-31T16:01:34.110501+00:00
> Content-Hash: db9ea4350650d346
> Type: html

---

### Overview ¶

Package unarr is a decompression library for RAR, TAR, ZIP and 7z archives. 

### Index ¶

  * Variables
  * type Archive
  *     * func NewArchive(path string) (a *Archive, err error)
    * func NewArchiveFromMemory(b []byte) (a *Archive, err error)
    * func NewArchiveFromReader(r io.Reader) (a *Archive, err error)
  *     * func (a *Archive) Close() (err error)
    * func (a *Archive) Entry() error
    * func (a *Archive) EntryAt(off int64) error
    * func (a *Archive) EntryFor(name string) error
    * func (a *Archive) Extract(path string) (contents []string, err error)
    * func (a *Archive) List() (contents []string, err error)
    * func (a *Archive) ModTime() time.Time
    * func (a *Archive) Name() string
    * func (a *Archive) Offset() int64
    * func (a *Archive) RawName() string
    * func (a *Archive) Read(b []byte) (n int, err error)
    * func (a *Archive) ReadAll() ([]byte, error)
    * func (a *Archive) Seek(offset int64, whence int) (int64, error)
    * func (a *Archive) Size() int



### Constants ¶

This section is empty.

### Variables ¶

[View Source](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L16)
    
    
    var (
    	ErrOpenFile    = [errors](/errors).[New](/errors#New)("unarr: open file failed")
    	ErrOpenMemory  = [errors](/errors).[New](/errors#New)("unarr: open memory failed")
    	ErrOpenArchive = [errors](/errors).[New](/errors#New)("unarr: no valid RAR, ZIP, 7Z or TAR archive")
    	ErrEntry       = [errors](/errors).[New](/errors#New)("unarr: failed to parse entry")
    	ErrEntryAt     = [errors](/errors).[New](/errors#New)("unarr: failed to parse entry at")
    	ErrEntryFor    = [errors](/errors).[New](/errors#New)("unarr: failed to parse entry for")
    	ErrSeek        = [errors](/errors).[New](/errors#New)("unarr: seek failed")
    	ErrRead        = [errors](/errors).[New](/errors#New)("unarr: read failure")
    )

### Functions ¶

This section is empty.

### Types ¶

####  type [Archive](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L28) ¶
    
    
    type Archive struct {
    	// contains filtered or unexported fields
    }

Archive represents unarr archive 

####  func [NewArchive](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L36) ¶
    
    
    func NewArchive(path [string](/builtin#string)) (a *Archive, err [error](/builtin#error))

NewArchive returns new unarr Archive 

####  func [NewArchiveFromMemory](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L51) ¶
    
    
    func NewArchiveFromMemory(b [][byte](/builtin#byte)) (a *Archive, err [error](/builtin#error))

NewArchiveFromMemory returns new unarr Archive from byte slice 

####  func [NewArchiveFromReader](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L66) ¶
    
    
    func NewArchiveFromReader(r [io](/io).[Reader](/io#Reader)) (a *Archive, err [error](/builtin#error))

NewArchiveFromReader returns new unarr Archive from io.Reader 

####  func (*Archive) [Close](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L164) ¶
    
    
    func (a *Archive) Close() (err [error](/builtin#error))

Close closes the underlying unarr archive 

####  func (*Archive) [Entry](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L103) ¶
    
    
    func (a *Archive) Entry() [error](/builtin#error)

Entry reads the next archive entry. 

io.EOF is returned when there is no more to be read from the archive. 

####  func (*Archive) [EntryAt](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L118) ¶
    
    
    func (a *Archive) EntryAt(off [int64](/builtin#int64)) [error](/builtin#error)

EntryAt reads the archive entry at the given offset 

####  func (*Archive) [EntryFor](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L128) ¶
    
    
    func (a *Archive) EntryFor(name [string](/builtin#string)) [error](/builtin#error)

EntryFor reads the (first) archive entry associated with the given name 

####  func (*Archive) [Extract](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L224) ¶
    
    
    func (a *Archive) Extract(path [string](/builtin#string)) (contents [][string](/builtin#string), err [error](/builtin#error))

Extract extracts archive to destination path 

####  func (*Archive) [List](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L258) ¶
    
    
    func (a *Archive) List() (contents [][string](/builtin#string), err [error](/builtin#error))

List lists the contents of archive 

####  func (*Archive) [ModTime](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L192) ¶
    
    
    func (a *Archive) ModTime() [time](/time).[Time](/time#Time)

ModTime returns the stored modification time of the current entry 

####  func (*Archive) [Name](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L182) ¶
    
    
    func (a *Archive) Name() [string](/builtin#string)

Name returns the name of the current entry as UTF-8 string 

####  func (*Archive) [Offset](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L177) ¶
    
    
    func (a *Archive) Offset() [int64](/builtin#int64)

Offset returns the stream offset of the current entry, for use with EntryAt 

####  func (*Archive) [RawName](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L187) ¶ added in v0.1.2
    
    
    func (a *Archive) RawName() [string](/builtin#string)

RawName returns the name of the current entry as raw string 

####  func (*Archive) [Read](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L140) ¶
    
    
    func (a *Archive) Read(b [][byte](/builtin#byte)) (n [int](/builtin#int), err [error](/builtin#error))

Read tries to read 'b' bytes into buffer, advancing the read offset pointer. 

Returns the actual number of bytes read. 

####  func (*Archive) [ReadAll](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L198) ¶
    
    
    func (a *Archive) ReadAll() ([][byte](/builtin#byte), [error](/builtin#error))

ReadAll reads current entry and returns data 

####  func (*Archive) [Seek](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L154) ¶
    
    
    func (a *Archive) Seek(offset [int64](/builtin#int64), whence [int](/builtin#int)) ([int64](/builtin#int64), [error](/builtin#error))

Seek moves the read offset pointer interpreted according to whence. 

Returns the new offset. 

####  func (*Archive) [Size](https://github.com/gen2brain/go-unarr/blob/v0.2.4/unarr.go#L172) ¶
    
    
    func (a *Archive) Size() [int](/builtin#int)

Size returns the total size of uncompressed data of the current entry 
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
