package traintracker

type TrainTrackerProps struct {
	Key string
}

type TrainTracker struct {
	key string
}

func NewTrainTracker(props TrainTrackerProps) TrainTracker {
	return TrainTracker{key: props.Key}
}
