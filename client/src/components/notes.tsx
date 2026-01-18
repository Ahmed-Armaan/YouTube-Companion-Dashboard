import { useEffect, useState } from "react"
import type { Note } from "../types/types"

type NotesProps = {
	videoId: string
}

type NotesResponse = {
	items: Note[]
	nextCursor?: string
}

function Notes({ videoId }: NotesProps) {
	const [notes, setNotes] = useState<Note[]>([])
	const [content, setContent] = useState<string>("")
	const [tags, setTags] = useState<string>("")

	const [loading, setLoading] = useState<boolean>(false)
	const [posting, setPosting] = useState<boolean>(false)

	// Fetch notes
	const fetchNotes = async () => {
		if (!videoId) return

		setLoading(true)
		try {
			const url = new URL(
				"/notes",
				import.meta.env.VITE_BACKEND_URL
			)
			url.searchParams.set("videoId", videoId)

			const res = await fetch(url.toString(), {
				credentials: "include",
			})

			if (!res.ok) throw new Error("Failed to fetch notes")

			const data: NotesResponse = await res.json()
			setNotes(data.items)
		} catch (err) {
			console.error(err)
		} finally {
			setLoading(false)
		}
	}

	// Create / update note
	const createNote = async () => {
		if (!content.trim() || posting) return

		setPosting(true)
		try {
			const res = await fetch(
				`${import.meta.env.VITE_BACKEND_URL}/notes`,
				{
					method: "POST",
					credentials: "include",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({
						videoId,
						content,
						tags: tags
							.split(",")
							.map(t => t.trim())
							.filter(Boolean),
					}),
				}
			)

			if (!res.ok) throw new Error("Failed to save note")

			const note: Note = await res.json()

			setNotes(prev => {
				const filtered = prev.filter(n => n.id !== note.id)
				return [note, ...filtered]
			})

			setContent("")
			setTags("")
		} catch (err) {
			console.error(err)
		} finally {
			setPosting(false)
		}
	}

	// Delete note
	const deleteNote = async (id: string) => {
		try {
			await fetch(
				`${import.meta.env.VITE_BACKEND_URL}/notes?id=${id}`,
				{
					method: "DELETE",
					credentials: "include",
				}
			)
			setNotes(prev => prev.filter(n => n.id !== id))
		} catch (err) {
			console.error(err)
		}
	}

	// Load on video change
	useEffect(() => {
		setNotes([])
		if (videoId) fetchNotes()
	}, [videoId])

	return (
		<div className="space-y-4">
			<h2 className="font-medium">Notes</h2>

			{/* New note */}
			<div className="space-y-2">
				<textarea
					className="w-full border rounded p-2 text-sm"
					placeholder="Write a note…"
					value={content}
					onChange={(e) => setContent(e.target.value)}
				/>

				<input
					className="w-full border rounded p-2 text-sm"
					placeholder="tags (comma separated)"
					value={tags}
					onChange={(e) => setTags(e.target.value)}
				/>

				<button
					onClick={createNote}
					disabled={posting}
					className="text-sm text-blue-600"
				>
					{posting ? "Saving…" : "Add note"}
				</button>
			</div>

			{/* States */}
			{loading && (
				<p className="text-sm text-gray-500">Loading notes…</p>
			)}

			{!loading && notes.length === 0 && (
				<p className="text-sm text-gray-500">No notes yet</p>
			)}

			{/* Notes list */}
			<ul className="space-y-2">
				{notes.map(note => (
					<li
						key={note.id}
						className="border rounded p-2 text-sm space-y-1"
					>
						<div>{note.content}</div>

						{note.tags?.length > 0 && (
							<div className="text-xs text-gray-500">
								{note.tags.map(t => `#${t}`).join(" ")}
							</div>
						)}

						<button
							onClick={() => deleteNote(note.id)}
							className="text-xs text-red-600"
						>
							Delete
						</button>
					</li>
				))}
			</ul>
		</div>
	)
}

export default Notes
