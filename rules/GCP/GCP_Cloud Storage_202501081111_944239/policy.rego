package bucket_with_uniform_access_control_disabled_4400004_142
import rego.v1

default risk := false
risk if {
    not isBucketLevelAccessEnabled
}

isBucketLevelAccessEnabled if {
    input.Bucket.iamConfiguration.uniformBucketLevelAccess.enabled == true
}