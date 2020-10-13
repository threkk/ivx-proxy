import { ChangeEvent, FormEvent, useEffect, useRef, useState } from 'react'

function checkValidURL(str: string): boolean {
  try {
    const url = new URL(str)
    const [_, extension] = url.pathname.split('.')
    return url.hostname.endsWith('ivoox.com') && extension === 'xml'
  } catch {
    return false
  }
}

export default function URLWidget() {
  const [target, setTarget] = useState(null)
  const [url, setURL] = useState(null)
  const [user, setUser] = useState('')
  const [pass, setPass] = useState('')

  const urlRef = useRef(null)
  const targetRef = useRef(null)

  useEffect(() => {
    urlRef.current.focus()
  }, [])

  useEffect(() => {
    targetRef.current.select()
  }, [target])

  function buildTarget(ev: FormEvent) {
    ev.preventDefault()
    const targetURL = new URL(window.location.toString())
    targetURL.pathname = '/'
    if (user !== '') targetURL.username = user
    if (pass !== '') targetURL.password = pass
    targetURL.search = `?feed=${encodeURI(url)}`
    setTarget(targetURL)
    return
  }

  return (
    <>
      <form onSubmit={buildTarget}>
        <div className='field'>
          <label className='label' htmlFor='name'>
            User (optional):
          </label>
          <div className='control'>
            <input
              className='input'
              type='text'
              id='name'
              placeholder='Your username (if any)'
              value={user}
              onChange={(e) => setUser(e.target.value)}
            />
          </div>
        </div>

        <div className='field'>
          <label className='label' htmlFor='pass'>
            Password (optional):
          </label>
          <div className='control'>
            <input
              className='input'
              type='password'
              id='pass'
              placeholder='Your password (if any)'
              value={pass}
              onChange={(e) => setPass(e.target.value)}
            />
          </div>
        </div>
        <div className='field'>
          <label className='label' htmlFor='original'>
            URL of the podcast:
          </label>
          <div className='control'>
            <input
              className={`input ${
                checkValidURL(url) ? 'is-success' : 'is-danger'
              }`}
              ref={urlRef}
              type='text'
              id='original'
              placeholder='Your iVoox Podcast url'
              onChange={(e) => setURL(e.target.value)}
            />
          </div>
          <>
            <p className='help is-danger'>
              {!checkValidURL(url) ? 'The URL must be a feed from iVoox.' : ''}
            </p>
          </>
        </div>
        <div className='field'>
          <label className='label' htmlFor='target'>
            Proxied URL:
          </label>
          <div className='control'>
            <input
              className='input'
              type='text'
              id='target'
              ref={targetRef}
              disabled={target === null}
              value={target}
              onChange={(e) => e.preventDefault()}
              readOnly
              placeholder='Paste this in your podcast application'
            />
          </div>
        </div>
        <div className='field'>
          <button
            className='button is-primary'
            disabled={checkValidURL(url) !== true}
          >
            Generate URL
          </button>
        </div>
      </form>
    </>
  )
}
