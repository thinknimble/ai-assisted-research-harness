Web Hooks
Event subscriptions enable notifications to be sent to a publicationUrl when an event occurs. The notification will also contain data specific to the triggering event. The events currently supported are enumerated below, along with the data accompanying the event notification.

controlpanel.activateSurvey

Attribute	Type	Description
Status	String	The result will be Active if successful
SurveyID	GUID	The Survey ID of the survey that has been activated
BrandID	String	The ID of the user's organization
controlpanel.deactivateSurvey

Attribute	Type	Description
Status	String	The result will be Inactive if successful
SurveyID	GUID	The Survey ID of the survey that has been deactivated
BrandID	String	The ID of the user's organization
surveyengine.startedRecipientSession.{SurveyID}

Attribute	Type	Description
Status	String	The result will be Started if successful
SurveyID	GUID	The Survey ID of the survey which has been started
RecipientID	GUID	The Recipient ID of the response that has been started
ResponseEventContext	String	Custom context that can be set via the ResponseEventContext embedded data. (50 character limit)
DistributionID	GUID	The ID of the distribution that sent the survey
StartedDate	String	The date and time (UTC) that the response was started
BrandID	String	The ID of the user's organization
surveyengine.partialResponse.{SurveyID}

Attribute	Type	Description
Status	String	The status will always be Partial if successful
SurveyID	GUID	The Survey ID of the survey whose partial response has been recorded
RecipientID	GUID	The Recipient ID of the response that has been completed
ResponseEventContext	String	Custom context that can be set via the ResponseEventContext embedded data. (50 character limit)
DistributionID	GUID	The ID of the distribution that sent the survey
ResponseID	GUID	The Response ID of the response that has been partially completed
CompletedDate	String	The date and time (UTC) that the response was recorded
BrandID	String	The ID of the user's organization
surveyengine.completedResponse.{SurveyID}

Attribute	Type	Description
Status	String	The status will always be Complete if successful
SurveyID	GUID	The Survey ID of the survey whose response has been completed
ResponseID	GUID	The Response ID of the response that has been completed
RecipientID	GUID	The Recipient ID of the response that has been completed
ResponseEventContext	String	Custom context that can be set via the ResponseEventContext embedded data. (50 character limit)
DistributionID	GUID	The ID of the distribution that sent the survey
CompletedDate	String	The date and time (UTC) that the response was completed
BrandID	String	The ID of the user's organization