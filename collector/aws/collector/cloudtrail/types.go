package cloudtrail

import "time"

type TrailStatus struct {
	IsLogging                          *bool
	LatestCloudWatchLogsDeliveryError  *string
	LatestCloudWatchLogsDeliveryTime   *time.Time
	LatestDeliveryAttemptSucceeded     *string
	LatestDeliveryAttemptTime          *string
	LatestDeliveryError                *string
	LatestDeliveryTime                 *time.Time
	LatestDigestDeliveryError          *string
	LatestDigestDeliveryTime           *time.Time
	LatestNotificationAttemptSucceeded *string
	LatestNotificationAttemptTime      *string
	LatestNotificationError            *string
	LatestNotificationTime             *time.Time
	StartLoggingTime                   *time.Time
	StopLoggingTime                    *time.Time
	TimeLoggingStarted                 *string
	TimeLoggingStopped                 *string
}
