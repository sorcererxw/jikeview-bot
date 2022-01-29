package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	/**
	{
	    "streams": [
	        {
	            "index": 0,
	            "codec_name": "h264",
	            "codec_long_name": "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
	            "profile": "High",
	            "codec_type": "video",
	            "codec_time_base": "794561/47673600",
	            "codec_tag_string": "avc1",
	            "codec_tag": "0x31637661",
	            "width": 1920,
	            "height": 1072,
	            "coded_width": 1920,
	            "coded_height": 1072,
	            "has_b_frames": 2,
	            "pix_fmt": "yuv420p",
	            "level": 40,
	            "chroma_location": "left",
	            "refs": 1,
	            "is_avc": "true",
	            "nal_length_size": "4",
	            "r_frame_rate": "30/1",
	            "avg_frame_rate": "23836800/794561",
	            "time_base": "1/16000",
	            "start_pts": 0,
	            "start_time": "0.000000",
	            "duration_ts": 3972805,
	            "duration": "248.300313",
	            "bit_rate": "2471935",
	            "bits_per_raw_sample": "8",
	            "nb_frames": "7449",
	            "disposition": {
	                "default": 1,
	                "dub": 0,
	                "original": 0,
	                "comment": 0,
	                "lyrics": 0,
	                "karaoke": 0,
	                "forced": 0,
	                "hearing_impaired": 0,
	                "visual_impaired": 0,
	                "clean_effects": 0,
	                "attached_pic": 0,
	                "timed_thumbnails": 0
	            },
	            "tags": {
	                "language": "und",
	                "handler_name": "VideoHandler"
	            }
	        },
	        {
	            "index": 1,
	            "codec_name": "aac",
	            "codec_long_name": "AAC (Advanced Audio Coding)",
	            "profile": "LC",
	            "codec_type": "audio",
	            "codec_time_base": "1/44100",
	            "codec_tag_string": "mp4a",
	            "codec_tag": "0x6134706d",
	            "sample_fmt": "fltp",
	            "sample_rate": "44100",
	            "channels": 2,
	            "channel_layout": "stereo",
	            "bits_per_sample": 0,
	            "r_frame_rate": "0/0",
	            "avg_frame_rate": "0/0",
	            "time_base": "1/44100",
	            "start_pts": 0,
	            "start_time": "0.000000",
	            "duration_ts": 10953699,
	            "duration": "248.383197",
	            "bit_rate": "128012",
	            "max_bit_rate": "128012",
	            "nb_frames": "10697",
	            "disposition": {
	                "default": 1,
	                "dub": 0,
	                "original": 0,
	                "comment": 0,
	                "lyrics": 0,
	                "karaoke": 0,
	                "forced": 0,
	                "hearing_impaired": 0,
	                "visual_impaired": 0,
	                "clean_effects": 0,
	                "attached_pic": 0,
	                "timed_thumbnails": 0
	            },
	            "tags": {
	                "language": "und",
	                "handler_name": "SoundHandler"
	            }
	        }
	    ],
	    "format": {
	        "filename": "12d10526-d79b-427a-ab83-07f371319b93.mp4",
	        "nb_streams": 2,
	        "nb_programs": 0,
	        "format_name": "mov,mp4,m4a,3gp,3g2,mj2",
	        "format_long_name": "QuickTime / MOV",
	        "start_time": "0.000000",
	        "duration": "248.384000",
	        "size": "81006484",
	        "bit_rate": "2609072",
	        "probe_score": 100,
	        "tags": {
	            "major_brand": "isom",
	            "minor_version": "512",
	            "compatible_brands": "isomiso2avc1mp41",
	            "encoder": "Lavf58.29.100",
	            "description": "Packed by Bilibili XCoder v2.0.2"
	        }
	    }
	}
	*/
	ProbeData struct {
		Streams []struct {
			Index            int    `json:"index"`
			CodecName        string `json:"codec_name"`
			CodecLongName    string `json:"codec_long_name"`
			Profile          string `json:"profile"`
			CodecType        string `json:"codec_type"`
			CodecTimeBase    string `json:"codec_time_base"`
			CodecTagString   string `json:"codec_tag_string"`
			CodecTag         string `json:"codec_tag"`
			Width            int    `json:"width,omitempty"`
			Height           int    `json:"height,omitempty"`
			CodedWidth       int    `json:"coded_width,omitempty"`
			CodedHeight      int    `json:"coded_height,omitempty"`
			HasBFrames       int    `json:"has_b_frames,omitempty"`
			PixFmt           string `json:"pix_fmt,omitempty"`
			Level            int    `json:"level,omitempty"`
			ChromaLocation   string `json:"chroma_location,omitempty"`
			Refs             int    `json:"refs,omitempty"`
			IsAvc            string `json:"is_avc,omitempty"`
			NalLengthSize    string `json:"nal_length_size,omitempty"`
			RFrameRate       string `json:"r_frame_rate"`
			AvgFrameRate     string `json:"avg_frame_rate"`
			TimeBase         string `json:"time_base"`
			StartPts         int    `json:"start_pts"`
			StartTime        string `json:"start_time"`
			DurationTs       int    `json:"duration_ts"`
			Duration         string `json:"duration"`
			BitRate          string `json:"bit_rate"`
			BitsPerRawSample string `json:"bits_per_raw_sample,omitempty"`
			NbFrames         string `json:"nb_frames"`
			Disposition      struct {
				Default         int `json:"default"`
				Dub             int `json:"dub"`
				Original        int `json:"original"`
				Comment         int `json:"comment"`
				Lyrics          int `json:"lyrics"`
				Karaoke         int `json:"karaoke"`
				Forced          int `json:"forced"`
				HearingImpaired int `json:"hearing_impaired"`
				VisualImpaired  int `json:"visual_impaired"`
				CleanEffects    int `json:"clean_effects"`
				AttachedPic     int `json:"attached_pic"`
				TimedThumbnails int `json:"timed_thumbnails"`
			} `json:"disposition"`
			Tags struct {
				Language    string `json:"language"`
				HandlerName string `json:"handler_name"`
			} `json:"tags"`
			SampleFmt     string `json:"sample_fmt,omitempty"`
			SampleRate    string `json:"sample_rate,omitempty"`
			Channels      int    `json:"channels,omitempty"`
			ChannelLayout string `json:"channel_layout,omitempty"`
			BitsPerSample int    `json:"bits_per_sample,omitempty"`
			MaxBitRate    string `json:"max_bit_rate,omitempty"`
		} `json:"streams"`
		Format struct {
			Filename       string `json:"filename"`
			NbStreams      int    `json:"nb_streams"`
			NbPrograms     int    `json:"nb_programs"`
			FormatName     string `json:"format_name"`
			FormatLongName string `json:"format_long_name"`
			StartTime      string `json:"start_time"`
			Duration       string `json:"duration"`
			Size           string `json:"size"`
			BitRate        string `json:"bit_rate"`
			ProbeScore     int    `json:"probe_score"`
			Tags           struct {
				MajorBrand       string `json:"major_brand"`
				MinorVersion     string `json:"minor_version"`
				CompatibleBrands string `json:"compatible_brands"`
				Encoder          string `json:"encoder"`
				Description      string `json:"description"`
			} `json:"tags"`
		} `json:"format"`
	}
)

func DownloadAndFormatAudio(url string) (string, error) {
	return "", nil
}

func DownloadAndFormatVideo(url string) (string, error) {
	outputFile := path.Join(os.TempDir(), uuid.New().String()+".mp4")
	cmd := exec.Command(
		"ffmpeg",
		"-i", url,
		"-c", "copy",
		"-loglevel", "panic",
		"-bsf:a", "aac_adtstoasc",
		outputFile,
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return outputFile, nil
}

func GetVideoInfo(file string) (*ProbeData, error) {
	cmd := exec.Command("ffprobe",
		"-i", file,
		"-print_format", "json",
		"-loglevel", "fatal",
		"-show_format", "-show_streams")
	var outputBuf bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &outputBuf
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if stdErr.String() != "" {
		return nil, fmt.Errorf("ffprobe err: %s", stdErr.String())
	}
	var data ProbeData
	err = json.Unmarshal(outputBuf.Bytes(), &data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &data, nil
}

func DownloadImage(url string) (string, error) {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(url)))
	res, err := http.Get(url)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer res.Body.Close()
	filepath := path.Join(os.TempDir(), hash)
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer file.Close()
	if _, err := io.Copy(file, res.Body); err != nil {
		return "", errors.WithStack(err)
	}
	return filepath, nil
}
