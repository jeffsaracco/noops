#
# Ruby automated fizzbot example
#
# Can you write a program that passes the fizzbot test?
#
require "net/http"
require "json"

def main
  start = get_json('/fizzbot')

  first_question_path = start['nextQuestion']

  get_json(first_question_path)

  answer_result = send_answer(first_question_path, 'Ruby')

  num = 1
  while answer_result['result'] == 'correct' do
    question_path = answer_result['nextQuestion']
    puts "Question #{num}"
    question = get_json(question_path)

    answer = get_answer(question, num)

    answer_result = send_answer(question_path, answer)
    num += 1
  end
end

def send_answer(path, answer)
  post_json(path, { :answer => answer })
end

# get data from the api and parse it into a ruby hash
def get_json(path)
  response = Net::HTTP.get_response(build_uri(path))
  result = JSON.parse(response.body)

  puts JSON.pretty_generate(result)
  result
end

# post an answer to the noops api
def post_json(path, body)
  uri = build_uri(path)

  post_request = Net::HTTP::Post.new(uri, 'Content-Type' => 'application/json')
  post_request.body = JSON.generate(body)

  response = Net::HTTP.start(uri.hostname, uri.port, :use_ssl => true) do |http|
    http.request(post_request)
  end

  puts "HTTP #{response.code}"
  result = JSON.parse(response.body)
  puts JSON.pretty_generate(result)
  result
end

def build_uri(path)
  URI.parse("https://api.noopschallenge.com" + path)
end

def get_answer(question, num)
  rules = get_rules(question)
  arr = question['numbers']
  arr.each_with_index do |a, i|
    if rules.any? { |n, r| a % n == 0 }
      if rules.all? { |n, r| a % n == 0 }
        arr[i] = rules.values.join('')
      else
        rules.each do |n, r|
          if a % n == 0
            arr[i] = r
            break
          end
        end
      end
    end
  end
  # puts arr.join(' ').split('')
  arr.join(' ')
end

def get_rules(question)
  h = {}
  question['rules'].each do |r|
    h[r['number']] = r['response']
  end

  h
end

main
