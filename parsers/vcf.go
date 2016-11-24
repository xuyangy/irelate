package parsers

import (
	"io"

	"github.com/xuyangy/irelate/interfaces"
	"github.com/xuyangy/vcfgo"
)

type Variant struct {
	interfaces.IVariant
	source  uint32
	related []interfaces.Relatable
}

func (v *Variant) String() string {
	return v.IVariant.String()
}

func NewVariant(v interfaces.IVariant, source uint32, related []interfaces.Relatable) *Variant {
	return &Variant{v, source, related}
}

func (v *Variant) AddRelated(r interfaces.Relatable) {
	if len(v.related) == 0 {
		v.related = make([]interfaces.Relatable, 0, 2)
	}
	v.related = append(v.related, r)
}

func (v *Variant) Related() []interfaces.Relatable {
	return v.related
}

func (v *Variant) SetSource(src uint32) { v.source = src }
func (v *Variant) Source() uint32       { return v.source }

func Vopen(rdr io.Reader, hdr *vcfgo.Header) (*vcfgo.Reader, error) {
	if hdr == nil {
		return vcfgo.NewReader(rdr, true)
	}
	return vcfgo.NewWithHeader(rdr, hdr, true)
}

func StreamVCF(vcf *vcfgo.Reader) interfaces.RelatableChannel {
	ch := make(interfaces.RelatableChannel, 256)
	go func() {
		j := 0
		for {
			v := vcf.Read()
			if v == nil {
				break
			}
			ch <- &Variant{v, 0, nil}
			j++
			if j < 1000 {
				vcf.Clear()
				j = 0
			}
		}
		close(ch)
	}()
	return ch
}

type vWrapper struct {
	*vcfgo.Reader
}

func (v vWrapper) Next() (interfaces.Relatable, error) {
	r := v.Read()
	if r == nil {
		return nil, io.EOF
	}
	return &Variant{r, 0, nil}, nil
}

func (v vWrapper) AddInfoToHeader(id, itype, number, description string) {
	v.Reader.AddInfoToHeader(id, itype, number, description)
}

func VCFIterator(buf io.Reader) (interfaces.RelatableIterator, *vcfgo.Reader, error) {
	v, err := Vopen(buf, nil)
	if err != nil {
		return nil, v, err
	}
	return vWrapper{v}, v, nil
}
