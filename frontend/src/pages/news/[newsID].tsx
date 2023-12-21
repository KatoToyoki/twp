import '@style/global.css';

import { Col, Row } from 'react-bootstrap';

import NotFound from '@components/NotFound';

import newsData from '@pages/home/newsData.json';
import { useParams } from 'react-router-dom';

interface Props {
  id: number | null;
  image_url: string;
  title: string;
  date: string;
  subTitle: string;
  content: string;
}

const EachNews = () => {
  const params = useParams();

  const data: Props = { id: null, image_url: '', title: '', date: '', subTitle: '', content: '' };
  const foundNews = newsData.find((news) => news.id.toString() === params.news_id);

  if (foundNews) {
    Object.assign(data, foundNews);

    return (
      <div style={{ padding: '10% 10% 0% 10%' }}>
        <div className='news_bg flex-wrapper'>
          <Row>
            <Col xs={12} md={4}>
              <img src={data.image_url} className='news_pic' />
            </Col>
            <Col xs={12} md={8}>
              <h4 className='inpage_title'>{data.title}</h4> <br />
              <span className='right'>{data.date}</span>
              <hr className='hr' />
              <p>{data.subTitle}</p>
              <p>{data.content}</p>
            </Col>
          </Row>
        </div>
      </div>
    );
  } else {
    return <NotFound />;
  }
};

export default EachNews;
