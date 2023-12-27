import '@style/global.css';

import { Col, Row } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFileUpload, faPen, faTrash } from '@fortawesome/free-solid-svg-icons';
import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';

import TButton from '@components/TButton';
import FormItem from '@components/FormItem';

interface TagProps {
  id: number;
  name: string;
}

interface ProductProps {
  name: string;
  description: string;
  price: number;
  image: File;
  expire_date: string;
  stock: number;
  enable: string;
  tags: number[];
}

interface RequestProps {
  name: string;
  seller_name: string;
}

const tagStyle = {
  borderRadius: '30px',
  background: ' var(--button_light)',
  padding: '1% 1% 1% 3%',
  color: 'white',
  margin: '5px 0 5px 5px',
  width: '100%',
};

const LeftBgStyle = {
  backgroundColor: 'rgba(255, 255, 255, 0.7)',
  boxShadow: '6px 4px 10px 2px rgba(0, 0, 0, 0.25)',
};

const EmptyGoods = () => {
  const navigate = useNavigate();

  const [tag, setTag] = useState('');
  const [tags, setTags] = useState<TagProps[]>([]);
  const [queryTags, setQueryTags] = useState<string[]>([]);
  const [tagExists, setTagExists] = useState(false);
  const [file, setFile] = useState<string | null>(null);

  const { register, setValue, handleSubmit } = useForm<ProductProps>({
    defaultValues: {
      name: 'new product',
      description: 'new product description',
      price: 0,
      image: undefined,
      expire_date: 'expire date',
      stock: 0,
      enable: 'true',
      tags: undefined,
    },
  });

  const addTag = useMutation({
    mutationFn: async (data: RequestProps) => {
      const response = await fetch('/api/seller/tag', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
        },
        body: JSON.stringify(data),
      });
      if (!response.ok) {
        throw new Error('add tag failed');
      }
      return response.json();
    },
    onSuccess: (responseData: TagProps) => {
      console.log('success!', responseData);
      setTags((prevTags) => [...prevTags, responseData]);
      setValue(
        'tags',
        tags.map((tag) => tag.id),
      );
    },
    onError: (error: Error) => {
      console.log('not right', error);
    },
  });

  const queryTag = useMutation({
    mutationFn: async (data: string) => {
      const response = await fetch(`/api/seller/tag?name=${data}`, {
        headers: {
          Accept: 'application/json',
        },
      });
      if (!response.ok) {
        throw new Error('query tag failed');
      }
      return response.json();
    },
    onSuccess: (responseData: TagProps[]) => {
      console.log('success to query', responseData);
      const tagNames = responseData.map((tag) => tag.name);
      setQueryTags(tagNames);
      setTagExists(tagNames.includes(tag));
    },
    onError: (error: Error) => {
      console.log('not right on query', error);
    },
  });

  const addProduct = useMutation({
    mutationFn: async (data: ProductProps) => {
      const response = await fetch('/api/seller/product', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
        },
        body: JSON.stringify(data),
      });
      if (!response.ok) {
        throw new Error('add tag failed');
      }
      return response.json();
    },

    onSuccess: (responseData: TagProps) => {
      console.log('success on add product', responseData);
    },
    onError: (error: Error) => {
      console.log('not right on add product', error);
    },
  });

  const addNewTag = async (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.keyCode === 229) return;

    if (event.key === 'Enter') {
      const input = event.currentTarget.value.trim();
      console.log(event.currentTarget.value);

      if (input == '') {
        Reset();
        return;
      }

      await queryTag.mutate(tag);
      if (tagExists) {
        alert('Tag already exists');
        Reset();
        return;
      }

      // TODO : seller name need to be change to corresponding user
      await addTag.mutate({ name: input, seller_name: 'user1' });

      Reset();
    }
  };

  const TagOnChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const tagName = e.target.value;
    setQueryTags([]);

    if (tagName == '') {
      Reset();
      return;
    }
    setTag(tagName);
    await queryTag.mutate(tagName);
  };

  const deleteTag = (index: number) => {
    setTags((prevTags) => prevTags.filter((_, i) => i !== index));
    console.log('delete', tags);
  };

  const Reset = () => {
    setTag('');
    setQueryTags([]);
    setTagExists(false);
  };

  const uploadFile = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      if (!e.target.files[0].name.match(/\.(jpg|jpeg|png|gif)$/i)) {
        alert('not an image');
      } else {
        setValue('image', e.target.files[0]);
        setFile(URL.createObjectURL(e.target.files[0]));
      }
    }
  };

  const onSubmit: SubmitHandler<ProductProps> = async (data) => {
    console.log(data);
    await addProduct.mutate(data);
    // navigate('/user/seller/manageProducts');
  };

  return (
    <div style={{ padding: '55px 12% 0 12%' }}>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Row>
          <Col xs={12} md={5} style={LeftBgStyle}>
            <div className='flex_wrapper' style={{ padding: '0 8% 10% 8%' }}>
              <div
                style={{
                  position: 'relative',
                  width: '100%',
                  height: '100%',
                  borderRadius: '0 0 30px 0',
                }}
              >
                <div
                  className='center'
                  style={{ backgroundColor: 'black', borderRadius: '0 0 30px 0' }}
                >
                  {file ? (
                    <div>
                      <img
                        src={file}
                        alt='File preview'
                        style={{ width: '100%', height: '100%', borderRadius: '0 0 30px 0' }}
                      />
                    </div>
                  ) : (
                    <div style={{ padding: '30% 5% 30% 5%' }}>
                      <FontAwesomeIcon icon={faFileUpload} size='6x' />
                    </div>
                  )}
                </div>
                <br />
                <Row
                  style={{
                    position: 'absolute',
                    zIndex: '1',
                    right: '0px',
                    bottom: '40px',
                  }}
                >
                  <Col xs={9}></Col>
                  <Col xs={3}>
                    <form method='post' encType='multipart/form-data'>
                      <label
                        htmlFor='file'
                        className='custom-file-upload'
                        style={{ minWidth: '40px' }}
                      >
                        <div className='button pointer center' style={{ padding: '10px' }}>
                          <FontAwesomeIcon icon={faPen} className='white_word' />
                        </div>
                      </label>

                      <input
                        id='file'
                        name='file'
                        type='file'
                        style={{ display: 'none' }}
                        onChange={uploadFile}
                      />
                    </form>
                  </Col>
                </Row>
              </div>
              <br />
              <span className='dark'>add more tags</span>

              <input
                type='text'
                placeholder='Input tags'
                className='quantity_box'
                value={tag}
                onChange={TagOnChange}
                onKeyDown={addNewTag}
                style={{ marginBottom: '10px' }}
                list='queryTags'
              />
              <datalist id='queryTags'>
                {queryTags.map((tag, index) => (
                  <option key={index} value={tag} />
                ))}
              </datalist>

              <Row xs='auto'>
                {tags.map((currentTag, index) => (
                  <Col style={tagStyle} key={index} className='center'>
                    <Row style={{ width: '100%' }} className='center'>
                      <Col xs={1} className='center'>
                        <FontAwesomeIcon
                          icon={faTrash}
                          className='white_word pointer'
                          onClick={() => deleteTag(index)}
                        />
                      </Col>
                      <Col xs={9} lg={10}>
                        <span style={{ wordBreak: 'break-all' }}>{currentTag.name}</span>
                      </Col>
                    </Row>
                  </Col>
                ))}
              </Row>

              <div style={{ height: '50px' }} />
              <TButton
                text='Delete Product'
                action={() => navigate('/user/seller/manageProducts')}
              />
              <TButton text='Confirm Changes' action={handleSubmit(onSubmit)} />
            </div>
          </Col>
          <Col xs={12} md={7}>
            <div style={{ padding: '7% 0% 7% 0%' }}>
              <FormItem label='Product Name'>
                <input type='text' {...register('name', { required: true })} />
              </FormItem>

              <FormItem label='Product Price'>
                <input type='text' {...register('price', { required: true })} />
              </FormItem>

              <FormItem label='Product Quantity'>
                <input type='text' {...register('stock', { required: true })} />
              </FormItem>

              <FormItem label='Product Introduction'>
                <textarea {...register('description', { required: true })} />
              </FormItem>

              <FormItem label='Best Before Date'>
                <textarea {...register('expire_date', { required: true })} />
              </FormItem>
            </div>
          </Col>
        </Row>
      </form>
    </div>
  );
};

export default EmptyGoods;
