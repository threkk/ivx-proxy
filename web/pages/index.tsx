import Head from 'next/head'
import URLWidget from '../components/url'

export default function Home(): JSX.Element {
  return (
    <div>
      <Head>
        <title>iVoox Proxy</title>
        <link rel='icon' href='/favicon.ico' />
      </Head>

      <main className='columns is-gapless'>
        <div className='column is-desktop is-one-fifth'></div>
        <div className='column is-desktop is-three-fifths'>
          <section className='section'>
            <h1 className='title'>iVoox Proxy</h1>
            <h2 className='subtitle'>
              Access iVoox Originals from the app of your choice.
            </h2>
            <div className='content'>
              <p>
                {
                  'This form will allow you to generate a proxied URL for your favourite iVoox podcast. Some rules:'
                }
              </p>
              <ul>
                <li>{'It only works with podcast from https://ivoox.com'}</li>
                <li>
                  {
                    'User name and password are optional. If you needed a user and password to access this page, you will likely need them ;)'
                  }
                </li>
                <li>
                  {
                    'Your credentials are not sent over the network, they are used locally to generate the URL.'
                  }
                </li>
              </ul>
            </div>
            <URLWidget />
          </section>
        </div>
      </main>
      <footer className='footer'>
        <div className='content has-text-centered'>
          <p>
            <small className='small'>
              Created by <a href='https://threkk.com'>threkk</a> to deal with
              the frustration of his favourite podcast moving to iVoox Originals
              only.
            </small>
          </p>
          <p>
            <small className='small'>
              You can find the source code and more information about the
              project on{' '}
              <a href='https://github.com/threkk/ivoox-proxy/'>Github</a>
            </small>
          </p>
        </div>
      </footer>
    </div>
  )
}
