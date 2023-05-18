package jpeg_parser

type JpegMarker = byte

const (
	jpegNUL   JpegMarker = '\x00'
	jpegTEM   JpegMarker = '\x01'
	jpegSIZ   JpegMarker = '\x51'
	jpegCOD   JpegMarker = '\x52'
	jpegCOC   JpegMarker = '\x53'
	jpegTLM   JpegMarker = '\x55'
	jpegPLM   JpegMarker = '\x57'
	jpegPLT   JpegMarker = '\x58'
	jpegQCD   JpegMarker = '\x5C'
	jpegQCC   JpegMarker = '\x5D'
	jpegRGN   JpegMarker = '\x5E'
	jpegPOC   JpegMarker = '\x5F'
	jpegPPM   JpegMarker = '\x60'
	jpegPPT   JpegMarker = '\x61'
	jpegCRG   JpegMarker = '\x63'
	jpegCOM   JpegMarker = '\x64'
	jpegSEC   JpegMarker = '\x65'
	jpegEPB   JpegMarker = '\x66'
	jpegESD   JpegMarker = '\x67'
	jpegEPC   JpegMarker = '\x68'
	jpegRED   JpegMarker = '\x69'
	jpegSOT   JpegMarker = '\x90'
	jpegSOP   JpegMarker = '\x91'
	jpegEPH   JpegMarker = '\x92'
	jpegSOD   JpegMarker = '\x93'
	jpegINSEC JpegMarker = '\x94'
	jpegSOF0  JpegMarker = '\xC0'
	jpegSOF1  JpegMarker = '\xC1'
	jpegSOF2  JpegMarker = '\xC2'
	jpegSOF3  JpegMarker = '\xC3'
	jpegDHT   JpegMarker = '\xC4'
	jpegSOF5  JpegMarker = '\xC5'
	jpegSOF6  JpegMarker = '\xC6'
	jpegSOF7  JpegMarker = '\xC7'
	jpegJPG   JpegMarker = '\xC8'
	jpegSOF9  JpegMarker = '\xC9'
	jpegSOF10 JpegMarker = '\xCA'
	jpegSOF11 JpegMarker = '\xCB'
	jpegDAC   JpegMarker = '\xCC'
	jpegSOF13 JpegMarker = '\xCD'
	jpegSOF14 JpegMarker = '\xCE'
	jpegSOF15 JpegMarker = '\xCF'
	jpegRST0  JpegMarker = '\xD0'
	jpegRST1  JpegMarker = '\xD1'
	jpegRST2  JpegMarker = '\xD2'
	jpegRST3  JpegMarker = '\xD3'
	jpegRST4  JpegMarker = '\xD4'
	jpegRST5  JpegMarker = '\xD5'
	jpegRST6  JpegMarker = '\xD6'
	jpegRST7  JpegMarker = '\xD7'
	jpegSOI   JpegMarker = '\xD8'
	jpegEOI   JpegMarker = '\xD9'
	jpegSOS   JpegMarker = '\xDA'
	jpegDQT   JpegMarker = '\xDB'
	jpegDNL   JpegMarker = '\xDC'
	jpegDRI   JpegMarker = '\xDD'
	jpegDHP   JpegMarker = '\xDE'
	jpegEXP   JpegMarker = '\xDF'
	jpegAPP0  JpegMarker = '\xE0'
	jpegAPP1  JpegMarker = '\xE1'
	jpegAPP2  JpegMarker = '\xE2'
	jpegAPP3  JpegMarker = '\xE3'
	jpegAPP4  JpegMarker = '\xE4'
	jpegAPP5  JpegMarker = '\xE5'
	jpegAPP6  JpegMarker = '\xE6'
	jpegAPP7  JpegMarker = '\xE7'
	jpegAPP8  JpegMarker = '\xE8'
	jpegAPP9  JpegMarker = '\xE9'
	jpegAPP10 JpegMarker = '\xEA'
	jpegAPP11 JpegMarker = '\xEB'
	jpegAPP12 JpegMarker = '\xEC'
	jpegAPP13 JpegMarker = '\xED'
	jpegAPP14 JpegMarker = '\xEE'
	jpegAPP15 JpegMarker = '\xEF'
	jpegJPG0  JpegMarker = '\xF0'
	jpegJPG1  JpegMarker = '\xF1'
	jpegJPG2  JpegMarker = '\xF2'
	jpegJPG3  JpegMarker = '\xF3'
	jpegJPG4  JpegMarker = '\xF4'
	jpegJPG5  JpegMarker = '\xF5'
	jpegJPG6  JpegMarker = '\xF6'
	jpegSOF48 JpegMarker = '\xF7'
	jpegLSE   JpegMarker = '\xF8'
	jpegJPG9  JpegMarker = '\xF9'
	jpegJPG10 JpegMarker = '\xFA'
	jpegJPG11 JpegMarker = '\xFB'
	jpegJPG12 JpegMarker = '\xFC'
	jpegJPG13 JpegMarker = '\xFD'
	jpegCOM_  JpegMarker = '\xFE'
)

var jpeg_marker_byte = map[string]byte{
	"NUL":    '\x00',
	"TEM":    '\x01',
	"SIZ":    '\x51',
	"COD":    '\x52',
	"COC":    '\x53',
	"TLM":    '\x55',
	"PLM":    '\x57',
	"PLT":    '\x58',
	"QCD":    '\x5C',
	"QCC":    '\x5D',
	"RGN":    '\x5E',
	"POC":    '\x5F',
	"PPM":    '\x60',
	"PPT":    '\x61',
	"CRG":    '\x63',
	"COM":    '\x64',
	"SEC":    '\x65',
	"EPB":    '\x66',
	"ESD":    '\x67',
	"EPC":    '\x68',
	"RED":    '\x69',
	"SOT":    '\x90',
	"SOP":    '\x91',
	"EPH":    '\x92',
	"SOD":    '\x93',
	"INSEC":  '\x94',
	"SOF0":   '\xC0',
	"SOF1":   '\xC1',
	"SOF2":   '\xC2',
	"SOF3":   '\xC3',
	"DHT":    '\xC4',
	"SOF5":   '\xC5',
	"SOF6":   '\xC6',
	"SOF7":   '\xC7',
	"JPG":    '\xC8',
	"SOF9":   '\xC9',
	"SOF10":  '\xCA',
	"SOF11":  '\xCB',
	"DAC":    '\xCC',
	"SOF13":  '\xCD',
	"SOF14":  '\xCE',
	"SOF15":  '\xCF',
	"RST0":   '\xD0',
	"RST1":   '\xD1',
	"RST2":   '\xD2',
	"RST3":   '\xD3',
	"RST4":   '\xD4',
	"RST5":   '\xD5',
	"RST6":   '\xD6',
	"RST7":   '\xD7',
	"SOI":    '\xD8',
	"EOI":    '\xD9',
	"SOS":    '\xDA',
	"DQT":    '\xDB',
	"DNL":    '\xDC',
	"DRI":    '\xDD',
	"DHP":    '\xDE',
	"EXP":    '\xDF',
	"APP0":   '\xE0',
	"APP1":   '\xE1',
	"APP2":   '\xE2',
	"APP3":   '\xE3',
	"APP4":   '\xE4',
	"APP5":   '\xE5',
	"APP6":   '\xE6',
	"APP7":   '\xE7',
	"APP8":   '\xE8',
	"APP9":   '\xE9',
	"APP10":  '\xEA',
	"APP11":  '\xEB',
	"APP12":  '\xEC',
	"APP13":  '\xED',
	"APP14":  '\xEE',
	"APP15":  '\xEF',
	"JPG0":   '\xF0',
	"JPG1":   '\xF1',
	"JPG2":   '\xF2',
	"JPG3":   '\xF3',
	"JPG4":   '\xF4',
	"JPG5":   '\xF5',
	"JPG6":   '\xF6',
	"SOF48":  '\xF7',
	"LSE":    '\xF8',
	"JPG9":   '\xF9',
	"JPG10":  '\xFA',
	"JPG11":  '\xFB',
	"JPG12":  '\xFC',
	"JPG13":  '\xFD',
	"COMEXT": '\xFE',
}
