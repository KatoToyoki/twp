import { Col, Row } from 'react-bootstrap';
import { useQuery } from '@tanstack/react-query';

import GoodsItem from '@components/GoodsItem';
import { Props } from '@components/GoodsItem';
import { CheckStatus } from '@lib/CheckStatus';

const Shop = () => {
  const { status, data } = useQuery({
    queryKey: ['shopsView'],
    queryFn: async () => {
      const response = await fetch(`/api/seller/product?offset=${0}&limit=${8}`, {
        headers: {
          'Content-Type': 'application/json',
        },
      });
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.json();
    },
  });

  if (status != 'success') {
    return <CheckStatus status={status} />;
  }

  return (
    <div>
      <div className='title'>All products</div>
      <hr className='hr' />
      <Row>
        {data.map((d: Props, index: number) => {
          return (
            <Col xs={6} md={3} key={index}>
              <GoodsItem id={d.id} name={d.name} image_url={d.image_url} />
            </Col>
          );
        })}
      </Row>
    </div>
  );
};

export default Shop;
