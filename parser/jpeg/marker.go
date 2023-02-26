package jpeg

type JpegMarker byte

const (
	NUL   JpegMarker = '\x00'
	TEM   JpegMarker = '\x01'
	SIZ   JpegMarker = '\x51'
	COD   JpegMarker = '\x52'
	COC   JpegMarker = '\x53'
	TLM   JpegMarker = '\x55'
	PLM   JpegMarker = '\x57'
	PLT   JpegMarker = '\x58'
	QCD   JpegMarker = '\x5C'
	QCC   JpegMarker = '\x5D'
	RGN   JpegMarker = '\x5E'
	POC   JpegMarker = '\x5F'
	PPM   JpegMarker = '\x60'
	PPT   JpegMarker = '\x61'
	CRG   JpegMarker = '\x63'
	COM   JpegMarker = '\x64'
	SEC   JpegMarker = '\x65'
	EPB   JpegMarker = '\x66'
	ESD   JpegMarker = '\x67'
	EPC   JpegMarker = '\x68'
	RED   JpegMarker = '\x69'
	SOT   JpegMarker = '\x90'
	SOP   JpegMarker = '\x91'
	EPH   JpegMarker = '\x92'
	SOD   JpegMarker = '\x93'
	INSEC JpegMarker = '\x94'
	SOF0  JpegMarker = '\xC0'
	SOF1  JpegMarker = '\xC1'
	SOF2  JpegMarker = '\xC2'
	SOF3  JpegMarker = '\xC3'
	DHT   JpegMarker = '\xC4'
	SOF5  JpegMarker = '\xC5'
	SOF6  JpegMarker = '\xC6'
	SOF7  JpegMarker = '\xC7'
	JPG   JpegMarker = '\xC8'
	SOF9  JpegMarker = '\xC9'
	SOF10 JpegMarker = '\xCA'
	SOF11 JpegMarker = '\xCB'
	DAC   JpegMarker = '\xCC'
	SOF13 JpegMarker = '\xCD'
	SOF14 JpegMarker = '\xCE'
	SOF15 JpegMarker = '\xCF'
	RST0  JpegMarker = '\xD0'
	RST1  JpegMarker = '\xD1'
	RST2  JpegMarker = '\xD2'
	RST3  JpegMarker = '\xD3'
	RST4  JpegMarker = '\xD4'
	RST5  JpegMarker = '\xD5'
	RST6  JpegMarker = '\xD6'
	RST7  JpegMarker = '\xD7'
	SOI   JpegMarker = '\xD8'
	EOI   JpegMarker = '\xD9'
	SOS   JpegMarker = '\xDA'
	DQT   JpegMarker = '\xDB'
	DNL   JpegMarker = '\xDC'
	DRI   JpegMarker = '\xDD'
	DHP   JpegMarker = '\xDE'
	EXP   JpegMarker = '\xDF'
	APP0  JpegMarker = '\xE0'
	APP1  JpegMarker = '\xE1'
	APP2  JpegMarker = '\xE2'
	APP3  JpegMarker = '\xE3'
	APP4  JpegMarker = '\xE4'
	APP5  JpegMarker = '\xE5'
	APP6  JpegMarker = '\xE6'
	APP7  JpegMarker = '\xE7'
	APP8  JpegMarker = '\xE8'
	APP9  JpegMarker = '\xE9'
	APP10 JpegMarker = '\xEA'
	APP11 JpegMarker = '\xEB'
	APP12 JpegMarker = '\xEC'
	APP13 JpegMarker = '\xED'
	APP14 JpegMarker = '\xEE'
	APP15 JpegMarker = '\xEF'
	JPG0  JpegMarker = '\xF0'
	JPG1  JpegMarker = '\xF1'
	JPG2  JpegMarker = '\xF2'
	JPG3  JpegMarker = '\xF3'
	JPG4  JpegMarker = '\xF4'
	JPG5  JpegMarker = '\xF5'
	JPG6  JpegMarker = '\xF6'
	SOF48 JpegMarker = '\xF7'
	LSE   JpegMarker = '\xF8'
	JPG9  JpegMarker = '\xF9'
	JPG10 JpegMarker = '\xFA'
	JPG11 JpegMarker = '\xFB'
	JPG12 JpegMarker = '\xFC'
	JPG13 JpegMarker = '\xFD'
	COM_  JpegMarker = '\xFE'
)

var JpegMarkers = map[string]byte{
	"NUL":   '\x00',
	"TEM":   '\x01',
	"SIZ":   '\x51',
	"COD":   '\x52',
	"COC":   '\x53',
	"TLM":   '\x55',
	"PLM":   '\x57',
	"PLT":   '\x58',
	"QCD":   '\x5C',
	"QCC":   '\x5D',
	"RGN":   '\x5E',
	"POC":   '\x5F',
	"PPM":   '\x60',
	"PPT":   '\x61',
	"CRG":   '\x63',
	"COM":   '\x64',
	"SEC":   '\x65',
	"EPB":   '\x66',
	"ESD":   '\x67',
	"EPC":   '\x68',
	"RED":   '\x69',
	"SOT":   '\x90',
	"SOP":   '\x91',
	"EPH":   '\x92',
	"SOD":   '\x93',
	"INSEC": '\x94',
	"SOF0":  '\xC0',
	"SOF1":  '\xC1',
	"SOF2":  '\xC2',
	"SOF3":  '\xC3',
	"DHT":   '\xC4',
	"SOF5":  '\xC5',
	"SOF6":  '\xC6',
	"SOF7":  '\xC7',
	"JPG":   '\xC8',
	"SOF9":  '\xC9',
	"SOF10": '\xCA',
	"SOF11": '\xCB',
	"DAC":   '\xCC',
	"SOF13": '\xCD',
	"SOF14": '\xCE',
	"SOF15": '\xCF',
	"RST0":  '\xD0',
	"RST1":  '\xD1',
	"RST2":  '\xD2',
	"RST3":  '\xD3',
	"RST4":  '\xD4',
	"RST5":  '\xD5',
	"RST6":  '\xD6',
	"RST7":  '\xD7',
	"SOI":   '\xD8',
	"EOI":   '\xD9',
	"SOS":   '\xDA',
	"DQT":   '\xDB',
	"DNL":   '\xDC',
	"DRI":   '\xDD',
	"DHP":   '\xDE',
	"EXP":   '\xDF',
	"APP0":  '\xE0',
	"APP1":  '\xE1',
	"APP2":  '\xE2',
	"APP3":  '\xE3',
	"APP4":  '\xE4',
	"APP5":  '\xE5',
	"APP6":  '\xE6',
	"APP7":  '\xE7',
	"APP8":  '\xE8',
	"APP9":  '\xE9',
	"APP10": '\xEA',
	"APP11": '\xEB',
	"APP12": '\xEC',
	"APP13": '\xED',
	"APP14": '\xEE',
	"APP15": '\xEF',
	"JPG0":  '\xF0',
	"JPG1":  '\xF1',
	"JPG2":  '\xF2',
	"JPG3":  '\xF3',
	"JPG4":  '\xF4',
	"JPG5":  '\xF5',
	"JPG6":  '\xF6',
	"SOF48": '\xF7',
	"LSE":   '\xF8',
	"JPG9":  '\xF9',
	"JPG10": '\xFA',
	"JPG11": '\xFB',
	"JPG12": '\xFC',
	"JPG13": '\xFD',
	"COM_":  '\xFE',
}
