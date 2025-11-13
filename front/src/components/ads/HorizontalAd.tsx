import GoogleAdSense from './GoogleAdSense';

const HorizontalAd = () => {
  return (
    <div className="my-6">
      <GoogleAdSense
        adSlot="0987654321"
        adFormat="horizontal"
        className="w-full h-[90px]"
        style={{ width: '100%', height: '90px' }}
      />
    </div>
  );
};

export default HorizontalAd;
