import { Link } from "react-router-dom";
import { Channel } from "@/lib/channels";
import { Play } from "lucide-react";
import { Card } from "@/components/ui/card";

interface ChannelCardProps {
  channel: Channel;
}

const ChannelCard = ({ channel }: ChannelCardProps) => {
  return (
    <Link to={`/channel/${channel.id}`} state={{ channelData: channel }}>
      <Card className="group relative overflow-hidden bg-card border-border hover:border-primary transition-all duration-300 hover:scale-105">
        <div className="aspect-video relative overflow-hidden">
          <img
            src={channel.logo}
            alt={channel.name}
            className="w-full h-full object-cover"
          />
          <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center">
            <Play className="h-12 w-12 text-primary" />
          </div>
        </div>
        <div className="p-4">
          <div className="flex items-start justify-between gap-2 mb-2">
            <h3 className="font-semibold text-foreground line-clamp-1">{channel.name}</h3>
            <span className="text-xs px-2 py-1 rounded-full bg-primary/20 text-primary whitespace-nowrap">
              {channel.category}
            </span>
          </div>
          <p className="text-sm text-muted-foreground line-clamp-2">{channel.description}</p>
          <div className="flex items-center gap-2 mt-3 text-xs text-muted-foreground">
            <span>{channel.language}</span>
            <span>•</span>
            <span>{channel.country}</span>
          </div>
        </div>
      </Card>
    </Link>
  );
};

export default ChannelCard;
