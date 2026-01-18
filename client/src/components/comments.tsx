import { useEffect, useState } from "react"
import type {
	CommentThreadItem,
	YTCommentThreadResponse,
} from "../types/types"

type CommentsProps = {
	videoId: string
}

function Comments({ videoId }: CommentsProps) {
	const [comments, setComments] = useState<CommentThreadItem[]>([])
	const [pageToken, setPageToken] = useState<string>("")
	const [hasMore, setHasMore] = useState(true)
	const [loading, setLoading] = useState(false)

	const [newComment, setNewComment] = useState("")
	const [replyingTo, setReplyingTo] = useState<string | null>(null)
	const [replyText, setReplyText] = useState("")

	const [myChannelId, setMyChannelId] = useState<string>("")

	useEffect(() => {
		fetch(`${import.meta.env.VITE_BACKEND_URL}/channelId`, {
			credentials: "include"
		})
			.then(r => r.json())
			.then(d => setMyChannelId(d.channelId))
	}, [])

	// fetch comments
	const fetchComments = async (reset = false) => {
		if (!videoId || loading || (!hasMore && !reset)) return

		setLoading(true)

		try {
			const backendURL = import.meta.env.VITE_BACKEND_URL
			const url = new URL("/comments", backendURL)

			url.searchParams.set("videoId", videoId)
			if (!reset && pageToken) {
				url.searchParams.set("pageToken", pageToken)
			}

			const res = await fetch(url.toString(), {
				credentials: "include",
			})

			if (!res.ok) throw new Error("Failed to fetch comments")

			const data: YTCommentThreadResponse = await res.json()

			setComments(prev =>
				reset ? data.items : [...prev, ...data.items]
			)
			setPageToken(data.nextPageToken || "")
			setHasMore(Boolean(data.nextPageToken))
		} catch (err) {
			console.error(err)
		} finally {
			setLoading(false)
		}
	}

	// post comments
	const postComment = async () => {
		if (!newComment.trim()) return

		try {
			const backendURL = import.meta.env.VITE_BACKEND_URL
			const url = new URL("/comments", backendURL)

			const res = await fetch(url.toString(), {
				method: "POST",
				credentials: "include",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					videoId,
					text: newComment,
				}),
			})

			if (!res.ok) throw new Error("Post failed")

			setNewComment("")
			setHasMore(true)
			fetchComments(true)
		} catch (err) {
			console.error(err)
		}
	}

	// reply
	const postReply = async (parentCommentId: string) => {
		if (!replyText.trim()) return

		try {
			const backendURL = import.meta.env.VITE_BACKEND_URL
			const url = new URL("/comments/reply", backendURL)

			await fetch(url.toString(), {
				method: "POST",
				credentials: "include",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					parentId: parentCommentId,
					text: replyText,
				}),
			})

			setReplyText("")
			setReplyingTo(null)
			fetchComments(true)
		} catch (err) {
			console.error(err)
		}
	}

	// delete
	const deleteComment = async (commentId: string) => {
		try {
			const backendURL = import.meta.env.VITE_BACKEND_URL
			const url = new URL("/comments", backendURL)
			url.searchParams.set("commentId", commentId)

			await fetch(url.toString(), {
				method: "DELETE",
				credentials: "include",
			})

			fetchComments(true)
		} catch (err) {
			console.error(err)
		}
	}

	useEffect(() => {
		setComments([])
		setPageToken("")
		setHasMore(true)
		fetchComments(true)
	}, [videoId])

	return (
		<div className="space-y-4">
			<h2 className="font-medium">Comments</h2>

			{/* Add comment */}
			<div className="space-y-2">
				<textarea
					className="w-full border p-2 text-sm"
					placeholder="Add a comment…"
					value={newComment}
					onChange={(e) => setNewComment(e.target.value)}
				/>
				<button
					onClick={postComment}
					className="text-sm text-blue-600"
				>
					Post
				</button>
			</div>

			{/* Threads */}
			{comments.map((thread) => {
				const top = thread.snippet.topLevelComment
				const topCommentId = top.id

				const isTopMine =
					top.snippet.authorChannelId.value === myChannelId

				return (
					<div key={thread.id} className="border rounded p-3 space-y-2">
						{/* Top-level comment */}
						<div className="flex gap-2">
							<img
								src={top.snippet.authorProfileImageUrl}
								className="w-8 h-8 rounded-full"
							/>

							<div className="flex-1">
								<div className="text-sm font-medium">
									{top.snippet.authorDisplayName}
								</div>

								<div className="text-sm">
									{top.snippet.textOriginal}
								</div>

								<div className="flex gap-3 text-xs mt-1">
									<button
										className="text-blue-600"
										onClick={() => setReplyingTo(topCommentId)}
									>
										Reply
									</button>

									{isTopMine && (
										<button
											className="text-red-600"
											onClick={() => deleteComment(topCommentId)}
										>
											Delete
										</button>
									)}
								</div>
							</div>
						</div>

						{/* Reply box */}
						{replyingTo === topCommentId && (
							<div className="pl-10 space-y-2">
								<textarea
									className="w-full border p-2 text-sm"
									placeholder="Write a reply…"
									value={replyText}
									onChange={(e) => setReplyText(e.target.value)}
								/>
								<button
									className="text-xs text-blue-600"
									onClick={() => postReply(topCommentId)}
								>
									Post reply
								</button>
							</div>
						)}

						{/* Replies */}
						{thread.replies?.comments?.length > 0 && (
							<div className="pl-10 space-y-2">
								{thread.replies.comments.map((reply) => {
									const isReplyMine =
										reply.snippet.authorChannelId.value === myChannelId

									return (
										<div key={reply.id} className="flex gap-2">
											<img
												src={reply.snippet.authorProfileImageUrl}
												className="w-6 h-6 rounded-full"
											/>

											<div className="flex-1">
												<div className="text-xs font-medium">
													{reply.snippet.authorDisplayName}
												</div>

												<div className="text-xs">
													{reply.snippet.textOriginal}
												</div>

												{isReplyMine && (
													<button
														className="text-xs text-red-600 mt-1"
														onClick={() => deleteComment(reply.id)}
													>
														Delete
													</button>
												)}
											</div>
										</div>
									)
								})}
							</div>
						)}
					</div>
				)
			})}

			{hasMore && (
				<button
					onClick={() => fetchComments()}
					disabled={loading}
					className="text-sm text-blue-600"
				>
					{loading ? "Loading…" : "Load more"}
				</button>
			)}
		</div>
	)
}

export default Comments
