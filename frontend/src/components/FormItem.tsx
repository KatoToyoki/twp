import { Row, Col } from 'react-bootstrap';

interface Props {
  label: string;
  children: JSX.Element;
  isCenter?: boolean;
}

const FormItemStyle = {
  margin: '3% 0% 3% 0%',
};

const FormItem = ({ label, children, isCenter = true }: Props) => {
  return (
    <Row style={FormItemStyle}>
      <Col xs={12} md={4} className='center_vertical'>
        {label}
      </Col>
      <Col xs={12} md={8} className={!isCenter ? '' : 'form_item_wrapper'}>
        {children}
      </Col>
    </Row>
  );
};

export default FormItem;
