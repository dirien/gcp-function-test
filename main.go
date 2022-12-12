package main

import (
	cloudfunction "github.com/pulumi/pulumi-google-native/sdk/go/google/cloudfunctions/v2"
	storage "github.com/pulumi/pulumi-google-native/sdk/go/google/storage/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a Google Cloud resource (Storage Bucket)
		bucket, err := storage.NewBucket(ctx, "bucket", nil)
		if err != nil {
			return err
		}

		bucketObject, err := storage.NewBucketObject(ctx, "zip", &storage.BucketObjectArgs{
			Bucket: bucket.Name,
			Source: pulumi.NewAssetArchive(map[string]interface{}{
				".": pulumi.NewFileArchive("./gofunc"),
			}),
		})
		if err != nil {
			return err
		}

		function, err := cloudfunction.NewFunction(ctx, "function", &cloudfunction.FunctionArgs{
			Environment: cloudfunction.FunctionEnvironmentGen1,
			BuildConfig: &cloudfunction.BuildConfigArgs{
				Source: &cloudfunction.SourceArgs{
					StorageSource: &cloudfunction.StorageSourceArgs{
						Bucket: bucket.Name,
						Object: bucketObject.Name,
					},
				},
				Runtime: pulumi.String("go119"),
			},
			Location: pulumi.String("europe-west3"),
			Labels: pulumi.StringMap{
				"foo": pulumi.String("bar"),
			},
		})
		if err != nil {
			return err
		}

		// Export the bucket self-link
		ctx.Export("function", function.Name)

		return nil
	})
}
