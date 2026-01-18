import { useLocation, useNavigate } from "react-router"
import { useEffect, useState } from "react"
import { FiEdit2, FiCheck, FiX } from "react-icons/fi"
import Header from "../components/header"
import Notes from "../components/notes"
import Comments from "../components/comments"
import type { VideoDetails } from "../types/types"

function VideoDashboard() {
	const location = useLocation()
	const navigate = useNavigate()

	const video = (location.state as { video?: VideoDetails } | null)?.video

	const [editingTitle, setEditingTitle] = useState<boolean>(false)
	const [editingDesc, setEditingDesc] = useState<boolean>(false)

	const [title, setTitle] = useState<string>("")
	const [description, setDescription] = useState<string>("")

	const [aiSuggestions, setAiSuggestions] = useState<string[]>([])
	const [aiLoading, setAiLoading] = useState<boolean>(false)

	useEffect(() => {
		if (!video) {
			navigate("/")
			return
		}
		setTitle(video.title)
		setDescription(video.description)
	}, [video, navigate])

	if (!video) return null

	/* ===== AI TITLE ===== */
	const aiImproveTitle = async () => {
		setAiLoading(true)
		setAiSuggestions([])

		try {
			const res = await fetch(
				`${import.meta.env.VITE_BACKEND_URL}/ai/title`,
				{
					method: "POST",
					credentials: "include",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({
						title,
						description,
					}),
				}
			)

			if (!res.ok) throw new Error("AI failed")

			const data: { suggestions: string[] } = await res.json()
			setAiSuggestions(data.suggestions ?? [])
		} catch (err) {
			console.error(err)
		} finally {
			setAiLoading(false)
		}
	}

	return (
		<>
			<Header />

			<div className="max-w-4xl mx-auto p-6 space-y-6">

				{/* ===== VIDEO ===== */}
				<div className="aspect-video bg-black">
					<iframe
						src={`https://www.youtube.com/embed/${video.videoId}`}
						className="w-full h-full"
						allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
						allowFullScreen
					/>
				</div>

				{/* ===== TITLE ===== */}
				<div className="space-y-2">
					<div className="flex items-center gap-2">
						<h2 className="font-medium">Title</h2>

						<button
							onClick={() => setEditingTitle(true)}
							className="text-sm text-blue-600 flex items-center gap-1"
						>
							<FiEdit2 size={14} /> Edit
						</button>

						<button
							onClick={aiImproveTitle}
							disabled={aiLoading}
							className="text-sm text-purple-600"
						>
							{aiLoading ? "Thinking…" : "AI"}
						</button>
					</div>

					{editingTitle ? (
						<div className="space-y-2">
							<input
								className="w-full border rounded p-2 text-sm"
								value={title}
								onChange={(e) => setTitle(e.target.value)}
							/>

							<div className="flex gap-2">
								<button
									onClick={async () => {
										await fetch(
											`${import.meta.env.VITE_BACKEND_URL}/video/title`,
											{
												method: "PUT",
												credentials: "include",
												headers: { "Content-Type": "application/json" },
												body: JSON.stringify({
													videoId: video.videoId,
													title,
												}),
											}
										)
										setEditingTitle(false)
									}}
									className="text-sm text-green-600 flex items-center gap-1"
								>
									<FiCheck size={14} /> Save
								</button>

								<button
									onClick={() => {
										setTitle(video.title)
										setEditingTitle(false)
									}}
									className="text-sm text-gray-500 flex items-center gap-1"
								>
									<FiX size={14} /> Cancel
								</button>
							</div>
						</div>
					) : (
						<h1 className="text-xl font-semibold">{title}</h1>
					)}

					{/* AI suggestions */}
					{aiSuggestions.length > 0 && (
						<div className="border rounded p-2 space-y-1">
							<div className="text-xs text-gray-500">AI suggestions</div>
							{aiSuggestions.map((s, i) => (
								<button
									key={i}
									onClick={() => {
										setTitle(s)
										setEditingTitle(true)
										setAiSuggestions([])
									}}
									className="block w-full text-left text-sm hover:bg-gray-100 p-1 rounded"
								>
									{s}
								</button>
							))}
						</div>
					)}
				</div>

				{/* ===== DESCRIPTION ===== */}
				<div className="space-y-2">
					<div className="flex items-center gap-2">
						<h2 className="font-medium">Description</h2>

						<button
							onClick={() => setEditingDesc(true)}
							className="text-sm text-blue-600 flex items-center gap-1"
						>
							<FiEdit2 size={14} /> Edit
						</button>
					</div>

					{editingDesc ? (
						<div className="space-y-2">
							<textarea
								className="w-full border rounded p-2 text-sm min-h-[120px]"
								value={description}
								onChange={(e) => setDescription(e.target.value)}
							/>

							<div className="flex gap-2">
								<button
									onClick={async () => {
										await fetch(
											`${import.meta.env.VITE_BACKEND_URL}/video/description`,
											{
												method: "PUT",
												credentials: "include",
												headers: { "Content-Type": "application/json" },
												body: JSON.stringify({
													videoId: video.videoId,
													description,
												}),
											}
										)
										setEditingDesc(false)
									}}
									className="text-sm text-green-600 flex items-center gap-1"
								>
									<FiCheck size={14} /> Save
								</button>

								<button
									onClick={() => {
										setDescription(video.description)
										setEditingDesc(false)
									}}
									className="text-sm text-gray-500 flex items-center gap-1"
								>
									<FiX size={14} /> Cancel
								</button>
							</div>
						</div>
					) : (
						<p className="text-sm whitespace-pre-line">{description}</p>
					)}
				</div>

				{/* ===== NOTES & COMMENTS ===== */}
				<Notes videoId={video.videoId} />
				<Comments videoId={video.videoId} />

			</div>
		</>
	)
}

export default VideoDashboard
