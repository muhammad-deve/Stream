import GoogleAdSense from './GoogleAdSense';

const VerticalAd = () => {
  return (
    <div className="sticky top-20 hidden xl:block">
      <GoogleAdSense
        adSlot="1234567890"
        adFormat="vertical"
        className="w-[160px] h-[600px]"
        style={{ width: '160px', height: '600px' }}
      />
    </div>
  );
};

export default VerticalAd;
