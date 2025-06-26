package object_public_access_4500001_140
import rego.v1

default risk := false
risk if {
    not isEnforcedPublicAccessPrevention
    bucketIamPolicyWithUnrestrictedPrincipal
}
risk if {
    not isEnforcedPublicAccessPrevention
    managedFolderIamPolicyWithUnrestrictedPrincipal
}
risk if {
    not isEnforcedPublicAccessPrevention
    not isBucketLevelAccessEnabled
    bucketAclWithUnrestrictedPrincipal
}
risk if {
    not isEnforcedPublicAccessPrevention
    not isBucketLevelAccessEnabled
    defaultObjectAclWithUnrestrictedPrincipal
}

isEnforcedPublicAccessPrevention if {
    input.Bucket.iamConfiguration.publicAccessPrevention == "enforced"
}
isBucketLevelAccessEnabled if {
    input.Bucket.iamConfiguration.uniformBucketLevelAccess.enabled == true
}

unrestrictedPrincipal := {"allUsers", "allAuthenticatedUsers"}
bucketIamPolicyWithUnrestrictedPrincipal if {
    some binding in input.IamPolicy.bindings
    unrestrictedPrincipal[binding.members[_]]
}

managedFolderIamPolicyWithUnrestrictedPrincipal if {
    some managedFolder in input.ManagedFolder
    some binding in managedFolder.IamPolicy.bindings
    unrestrictedPrincipal[binding.members[_]]
}

bucketAclWithUnrestrictedPrincipal if {
    some acl in input.Bucket.acl
    unrestrictedPrincipal[acl.entity]
}

defaultObjectAclWithUnrestrictedPrincipal if {
    some acl in input.Bucket.defaultObjectAcl
    unrestrictedPrincipal[acl.entity]
}
