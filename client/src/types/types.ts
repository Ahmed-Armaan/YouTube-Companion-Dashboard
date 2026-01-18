export type VideoList = {
	nextPageToken: string
	videos: VideoDetails[]
}

export type VideoDetails = {
	videoId: string
	title: string
	description: string
	publishedAt: string
	thumbnail: string
	duration: string
	viewCount: string
	likeCount: string
	dislikeCount: string
	embeddedhtml: string
}

// Comment thread related
export type CommentSnippet = {
	authorDisplayName: string
	authorProfileImageUrl: string
	authorChannelUrl: string
	authorChannelId: {
		value: string
	}
	textOriginal: string
}

export type TopLevelComment = {
	id: string
	snippet: CommentSnippet
}

export type ReplyComment = {
	id: string
	snippet: CommentSnippet
}

export type CommentThreadItem = {
	id: string
	snippet: {
		channelId: string
		topLevelComment: TopLevelComment
	}
	replies: {
		comments: ReplyComment[]
	}
}

export type YTCommentThreadResponse = {
	nextPageToken: string
	items: CommentThreadItem[]
}

// Notes
export type Note = {
	id: string
	content: string
	tags: string[]
	createdAt: string
}
