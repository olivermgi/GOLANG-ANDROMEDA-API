package vod

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	"cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
	"github.com/olivermgi/golang-andromeda-api/config"
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*120)
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

func (c *clientUploader) MoveFile(object string, dstName string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*120)
	defer cancel()

	src := c.cl.Bucket(c.bucketName).Object(object)
	dst := c.cl.Bucket(c.bucketName).Object(dstName)

	dst = dst.If(storage.Conditions{DoesNotExist: true})

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %w", dstName, object, err)
	}
	if err := src.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", object, err)
	}

	return nil
}

func (c *clientUploader) MoveFolder(sourcePrefix, destinationPrefix string) error {
	ctx := context.Background()
	it := c.cl.Bucket(c.bucketName).Objects(ctx, &storage.Query{Prefix: sourcePrefix})

	for {
		objAttrs, err := it.Next()
		if err == storage.ErrObjectNotExist {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list objects: %w", err)
		}

		sourceName := objAttrs.Name
		destinationName := strings.Replace(sourceName, sourcePrefix, destinationPrefix, 1)

		if err := c.MoveFile(sourceName, destinationName); err != nil {
			return fmt.Errorf("failed to move file: %w", err)
		}
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*120)
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*120)
	defer cancel()

	response, err := c.cl.GetJob(ctx, req)

	if err != nil {
		return "", fmt.Errorf("GetJob Error: %v", err)
	}

	return fmt.Sprintf("%v", response.State), nil
}

func (c *clientTranscoder) DeleteJob(jobID string) error {
	req := &transcoderpb.DeleteJobRequest{
		Name: jobID,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*120)
	defer cancel()

	err := c.cl.DeleteJob(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteJob: %w", err)
	}

	return nil
}
