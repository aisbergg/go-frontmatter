package copy

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	adrg "github.com/adrg/frontmatter"
	aisbergg "github.com/aisbergg/go-frontmatter/pkg/frontmatter"
)

var txt = `+++
{
    "menu": {
        "header": "SVG Viewer",
        "items": [
            {
                "id": "Open"
            },
            {
                "id": "OpenNew",
                "label": "Open New"
            },
            null
        ]
    }
}
+++
There is immense joy in just watching - watching all the little creatures in nature. The first step to doing anything is to believe you can do it. See it finished in your mind before you ever start. That's a son of a gun of a cloud. We'll throw some old gray clouds in here just sneaking around and having fun.

Maybe there's a little something happening right here. We'll make some happy little bushes here. In your world you can create anything you desire.

Decide where your cloud lives. Maybe he lives right in here. Let's build some happy little clouds up here. That's what painting is all about. It should make you feel good when you paint.

If it's not what you want - stop and change it. Don't just keep going and expect it will get better. Anyone can paint. It's a very cold picture, I may have to go get my coat. It’s about to freeze me to death. Let's get wild today.

Mix your color marbly don't mix it dead. Let's make a nice big leafy tree. That's crazy. We'll lay all these little funky little things in there. I'm a water fanatic. I love water. Let's put some highlights on these little trees. The sun wouldn't forget them.

See. We take the corner of the brush and let it play back-and-forth. Trees live in your fan brush, but you have to scare them out. In painting, you have unlimited power. You have the ability to move mountains. You've got to learn to fight the temptation to resist these things. Just let them happen. Nature is so fantastic, enjoy it. Let it make you happy.

Put your feelings into it, your heart, it's your world. How do you make a round circle with a square knife? That's your challenge for the day. Let's make a happy little mountain now. But we're not there yet, so we don't need to worry about it. Making all those little fluffies that live in the clouds. But they're very easily killed. Clouds are delicate.

--- Bob Ross
`

// dummyUnmarshal returns the length of frontmatter data, because we don't
// actually care about the unmarshaled data and we don't want to benchmark the
// unmarshalling process.
func dummyUnmarshal(data []byte, v interface{}) error {
	x := v.(*interface{})
	*x = len(data)
	return nil
}

func BenchmarkAisberggFrontmatter(b *testing.B) {
	formats := []*aisbergg.Format{
		aisbergg.NewFormat("+++", "+++", dummyUnmarshal),
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var data interface{}
		body, err := aisbergg.Parse(strings.NewReader(txt), &data, formats...)
		checkResult(b, body, data, err)
	}
}

type buffer struct {
	data []byte
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return &buffer{data: make([]byte, 0, 4096)}
	},
}

func BenchmarkAisberggFrontmatterWithBuffer(b *testing.B) {
	formats := []*aisbergg.Format{
		aisbergg.NewFormat("+++", "+++", dummyUnmarshal),
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var data interface{}
		buf := bufPool.Get().(*buffer)
		body, err := aisbergg.ParseWithBuffer(buf.data, strings.NewReader(txt), &data, formats...)
		checkResult(b, body, data, err)
		bufPool.Put(buf)
	}
}

var bufPool2 = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

func BenchmarkAdrgFrontmatter(b *testing.B) {
	formats := []*adrg.Format{
		adrg.NewFormat("+++", "+++", dummyUnmarshal),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var data interface{}
		body, err := adrg.Parse(strings.NewReader(txt), &data, formats...)
		checkResult(b, body, data, err)
	}
}

// checkResult checks that the result is valid.
func checkResult(b *testing.B, body []byte, data interface{}, err error) {
	if err != nil {
		b.Fatal(err.Error())
	}
	if len(body) == 0 {
		b.Fatal("invalid body")
	}
	if data == nil || len(body) != 1870 {
		b.Fatalf("invalid data, expected length of %d got %d", 1870, len(body))
	}
}
