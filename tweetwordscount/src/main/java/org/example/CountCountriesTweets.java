package org.example;

import org.apache.hadoop.io.Text;

import java.io.IOException;
import java.util.Optional;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;
import org.apache.hadoop.util.GenericOptionsParser;
import org.json.simple.JSONObject;
import org.json.simple.JSONValue;

public class CountCountriesTweets {
    
    public static class TokenizerMapper
            extends Mapper<Object, Text, Text, IntWritable>{

        private final static IntWritable one = new IntWritable(1);
        private Text word = new Text();
        private final String NOT_FOUND = "NOT_FOUND";

        //use getLocation or getCountryuCode depending of the field you want to extract the user location
        public void map(Object key, Text value, Context context
        ) throws IOException, InterruptedException {
            String line = value.toString();
            String countryCode = getCountryCode(line);
            word.set(countryCode);
            context.write(word, one);
        }

        private String getLocation(String jsonText) {
            try {
                JSONObject obj = (JSONObject) JSONValue.parse(jsonText);
                Optional<JSONObject> user = Optional.of((JSONObject) obj.get("user"));
                if(!user.isPresent()) {
                    return NOT_FOUND;
                }
                return Optional.of((String) user.get().get("location")).orElse(NOT_FOUND);
            } catch (Exception e){
                return NOT_FOUND;
            }
        }

        private String getCountryCode(String jsonText) {
            try {
                JSONObject obj = (JSONObject) JSONValue.parse(jsonText);
                Optional<JSONObject> status = Optional.of((JSONObject) obj.get("retweeted_status"));
                if (!status.isPresent()) {
                    return NOT_FOUND;
                }

                Optional<JSONObject> place = Optional.of((JSONObject) status.get().get("place"));
                if(!place.isPresent()) {
                    return NOT_FOUND;
                }
                
                Optional<String> countryCode = Optional.of((String) place.get().get("country_code"));

                return countryCode.orElse(NOT_FOUND);
            } catch (Exception e) {
                return NOT_FOUND;
            }
        }
    }

    public static class IntSumReducer
            extends Reducer<Text,IntWritable,Text,IntWritable> {
        private IntWritable result = new IntWritable();

        public void reduce(Text key, Iterable<IntWritable> values,
                           Context context
        ) throws IOException, InterruptedException {
            int sum = 0;
            for (IntWritable val : values) {
                sum += val.get();
            }
            result.set(sum);
            context.write(key, result);
        }
    }

    public static void main(String[] args) throws Exception {
        Configuration conf = new Configuration();
        String[] otherArgs = new GenericOptionsParser(conf, args).getRemainingArgs();
        if (otherArgs.length != 2) {
            System.err.println("Usage: wordcount <in> <out>");
            System.exit(2);
        }
        Job job = new Job(conf, "word count");
        job.setJarByClass(CountCountriesTweets.class);
        job.setMapperClass(TokenizerMapper.class);
        job.setCombinerClass(IntSumReducer.class);
        job.setReducerClass(IntSumReducer.class);
        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(IntWritable.class);
        FileInputFormat.addInputPath(job, new Path(otherArgs[0]));
        FileOutputFormat.setOutputPath(job, new Path(otherArgs[1]));
        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
