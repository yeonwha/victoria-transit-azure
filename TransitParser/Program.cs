using Google.Protobuf;
using TransitRealtime;
using ProtoBuf;

// Load stop names from CSV
var stopNames = new Dictionary<string, string>();
foreach (var line in File.ReadLines("Victoria_Regional_Transit_System_stops.csv").Skip(1))
{
    var parts = line.Split(',');
    if (parts.Length >= 3)
        stopNames[parts[0]] = parts[2];
}

// Read and parse .pb file
var data = File.ReadAllBytes("tripupdates.pb");
var feed = Serializer.Deserialize<FeedMessage>(new MemoryStream(data));

Console.WriteLine("=== Victoria Transit Real-time Data ===");
foreach (var entity in feed.Entities)
{
    if (entity.TripUpdate != null)
    {
        foreach (var stopUpdate in entity.TripUpdate.StopTimeUpdates)
        {
            var stopId = stopUpdate.StopId;
            var name = stopNames.GetValueOrDefault(stopId, "Unknown Stop");

            var delay = stopUpdate.Departure?.Delay ?? stopUpdate.Arrival?.Delay ?? 0;

            if (delay >= 60)
                Console.WriteLine($"Stop: {name} ({stopId}) | Delay: {delay} sec(s)");
        }
    }
}