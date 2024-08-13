package vod

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	"cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
	"github.com/olivermgi/golang-crud-practice/config"
)

type clientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

type clientTranscoder struct {
	cl        *transcoder.Client
	projectID string
	location  string
	inputURI  string
	outputURI string
}

var gcpConfig map[string]map[string]string
var Uploader *clientUploader
var Transcoder *clientTranscoder

func init() {
	gcpConfig = config.GetGcpConfig()

	uploadClient, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}

	Uploader = &clientUploader{
		cl:         uploadClient,
		projectID:  gcpConfig["common"]["project_id"],
		bucketName: gcpConfig["storage"]["bucketName"],
		uploadPath: "",
	}

	transcodeClient, err := transcoder.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create transcode client: %v", err)
	}
	// defer transcodeClient.Close()

	Transcoder = &clientTranscoder{
		cl:        transcodeClient,
		projectID: gcpConfig["common"]["project_id"],
		location:  gcpConfig["transcoder"]["location"],
		inputURI:  gcpConfig["transcoder"]["input_uri"],
		outputURI: gcpConfig["transcoder"]["output_uri"],
	}
}

func (c *clientUploader) UploadFile(file multipart.File, object string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func (c *clientTranscoder) TransformVideoFile(inputPath string, outputPaht string) (string, error) {
	inputURI := gcpConfig["transcoder"]["input_uri"] + inputPath    // "/2/28af70fc.mp4"
	outputURI := gcpConfig["transcoder"]["output_uri"] + outputPaht // "/2/28af70fc/"

	req := &transcoderpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", c.projectID, c.location),
		Job: &transcoderpb.Job{
			InputUri:  inputURI,
			OutputUri: outputURI,
			// JobConfig: &transcoderpb.Job_TemplateId{
			// 	TemplateId: preset,
			// },
		},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	response, err := c.cl.CreateJob(ctx, req)
	if err != nil {
		return "", fmt.Errorf("CreateJob Error: %v", err)
	}

	// fmt.Printf("Job: %v\n")
	return response.GetName(), nil
}

func (c *clientTranscoder) GetJobState(jobID string) (string, error) {
	req := &transcoderpb.GetJobRequest{
		Name: jobID,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	response, err := c.cl.GetJob(ctx, req)

	if err != nil {
		return "", fmt.Errorf("GetJob Error: %v", err)
	}

	return fmt.Sprintf("%v", response.State), nil
}
