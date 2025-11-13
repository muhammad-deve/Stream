import GoogleAdSense from './GoogleAdSense';

const SidebarAd = () => {
  return (
    <div className="w-[300px] h-[250px]">
      <GoogleAdSense
        adSlot="5678901234"
        adFormat="rectangle"
        className="w-full h-full"
        style={{ width: '300px', height: '250px' }}
      />
    </div>
  );
};

export default SidebarAd;
