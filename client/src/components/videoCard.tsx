import { useNavigate } from 'react-router'
import '../index.css'
import type { VideoDetails } from '../types/types'

type VideoCardProps = {
	video: VideoDetails
}

function VideoCard({ video }: VideoCardProps) {
	const navigate = useNavigate()
	const navigateToDashboard = () => {
		navigate("/video", {
			state: { video }
		})
	}

	return (
		<>
			<div className='border border-black rounded-lg m-5 p-5 flex flex-row' onClick={navigateToDashboard}>
				{/* thumbnail */}
				<div>
					<img src={video.thumbnail} className='w-40 h-24' />
				</div>

				{/* video details */}
				<div className='flex flex-col px-5'>
					<div className='text-xl'>
						{video.title}
					</div>
					<div>
						{video.description}
					</div>
					<div className='flex p-2'>
						<div>{video.publishedAt}</div>
						<div>{video.viewCount}</div>
					</div>
				</div>
			</div >
		</>
	)
}

export default VideoCard
