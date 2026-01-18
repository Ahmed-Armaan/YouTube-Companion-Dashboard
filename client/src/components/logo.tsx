function Logo(props: { subtitle: boolean }) {
  return (
    <>
      <div className="flex flex-col">
        <div className="font-semibold text-[1em]">
          Youtube Dashboard
        </div>
        {props.subtitle &&
          <div className="text-[0.55em] text-text/60">
            Manage your YouTube channel at one location
          </div>
        }
      </div>
    </>
  )
}

export default Logo
