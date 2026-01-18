import { useEffect, useState } from "react"
import VideoCard from "../components/videoCard"
import Header from "../components/header"
import type { VideoList, VideoDetails } from "../types/types"

function Home() {
	const [pageToken, updatePageToken] = useState("")
	const [videos, setVideos] = useState<VideoDetails[]>([])

	useEffect(() => {
		FetchVideoList()
		console.log(videos)
	}, [])

	const appendNewVideos = (newVideos: VideoDetails[]) => {
		setVideos((prevVideos) => [...prevVideos, ...newVideos])
	}

	const FetchVideoList = async () => {
		const reqUrl = new URL(import.meta.env.VITE_BACKEND_URL + "/channel")
		reqUrl.searchParams.set("pageToken", pageToken)
		console.log(reqUrl)

		const res = await fetch(reqUrl, {
			credentials: "include",
		})

		if (!res.ok) {
			console.log("Failed to fetch Videos")
		} else {
			console.log(res)
			const videoList: VideoList = await res.json()
			updatePageToken(videoList.nextPageToken)
			appendNewVideos(videoList.videos)
			console.log(videoList)
		}
	}

	return (
		<>
			<Header />
			<div className="p-5">
				{
					videos.map((video) => (
						<VideoCard video={video} />
					))
				}
			</div>
		</>
	)
}

export default Home 
